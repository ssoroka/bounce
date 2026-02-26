package levels

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CubeBoundary struct {
	Point
	Size
	Rotation    float32
	StrokeWidth float32
	Color       color.Color
	tl          Point
	tr          Point
	bl          Point
	br          Point
}

func NewCubeBoundary(x, y, w, h, strokeWidth float32, color color.Color) *CubeBoundary {
	c := &CubeBoundary{
		Point:       Point{X: x, Y: y},
		Size:        Size{W: w, H: h},
		StrokeWidth: strokeWidth,
		Color:       color,
	}
	c.RecalculateCorners()
	return c
}

func (b *CubeBoundary) Scale(factor float32) {
	// Scale the size, but keep the center of the cube the same
	centerX := b.X + b.W/2
	centerY := b.Y + b.H/2
	b.X = centerX - (b.W*factor)/2
	b.Y = centerY - (b.H*factor)/2
	b.Size = b.Size.Scale(factor)
	b.RecalculateCorners()
}

func (b *CubeBoundary) RecalculateCorners() {
	// Center of rotation
	cx := b.X + b.W/2
	cy := b.Y + b.H/2

	// Unrotated corners relative to center
	corners := []Point{
		{-b.W / 2, -b.H / 2}, // tl
		{b.W / 2, -b.H / 2},  // tr
		{-b.W / 2, b.H / 2},  // bl
		{b.W / 2, b.H / 2},   // br
	}

	sin, cos := float32(math.Sin(float64(b.Rotation))), float32(math.Cos(float64(b.Rotation)))

	rotate := func(p Point) Point {
		return Point{
			X: cx + p.X*cos - p.Y*sin,
			Y: cy + p.X*sin + p.Y*cos,
		}
	}

	b.tl = rotate(corners[0])
	b.tr = rotate(corners[1])
	b.bl = rotate(corners[2])
	b.br = rotate(corners[3])
}

func (b *CubeBoundary) Draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, b.tl.X, b.tl.Y, b.tr.X, b.tr.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.tl.X, b.tl.Y, b.bl.X, b.bl.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.bl.X, b.bl.Y, b.br.X, b.br.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.tr.X, b.tr.Y, b.br.X, b.br.Y, b.StrokeWidth, b.Color, true)

	// draw normals
	if debug {
		lLine, lNorm := normal(b.tl, b.bl)
		bLine, bNorm := normal(b.bl, b.br)
		rLine, rNorm := normal(b.br, b.tr)
		tLine, tNorm := normal(b.tr, b.tl)

		vector.StrokeLine(screen, lLine.From.X, lLine.From.Y, lLine.To.X, lLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, bLine.From.X, bLine.From.Y, bLine.To.X, bLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, rLine.From.X, rLine.From.Y, rLine.To.X, rLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, tLine.From.X, tLine.From.Y, tLine.To.X, tLine.To.Y, b.StrokeWidth, green, true)
		// dot(lNorm.X, 0, lNorm.Y, 0)
		_ = lNorm
		_ = bNorm
		_ = rNorm
		_ = tNorm

	}
}

func (b *CubeBoundary) Update(delta float32) error {
	return nil
}

// CheckCircleCollision checks if a circle trapped within a cube boundary is colliding with the
// boundary walls and returns collision info to keep it inside.
func (b *CubeBoundary) CheckCircleCollision(c *Circle) Collision {
	lLine, lNorm := normal(b.tl, b.bl)
	bLine, bNorm := normal(b.bl, b.br)
	rLine, rNorm := normal(b.br, b.tr)
	tLine, tNorm := normal(b.tr, b.tl)

	// Helper to check if moving toward the wall
	movingToward := func(c *Circle, norm Vector) bool {
		return dot(c.Velocity.X, norm.X, c.Velocity.Y, norm.Y) < 0
	}

	// check left wall
	leftDist := dot(lNorm.X, c.X-lLine.From.X, lNorm.Y, c.Y-lLine.From.Y)
	if leftDist < c.Radius && movingToward(c, lNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: lNorm,
			Depth:  c.Radius - leftDist,
			Point:  Vector{X: c.X - lNorm.X*leftDist, Y: c.Y - lNorm.Y*leftDist},
		}
	}

	// check right wall
	rightDist := dot(rNorm.X, c.X-rLine.From.X, rNorm.Y, c.Y-rLine.From.Y)
	if rightDist < c.Radius && movingToward(c, rNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: rNorm,
			Depth:  c.Radius - rightDist,
			Point:  Vector{X: c.X - rNorm.X*rightDist, Y: c.Y - rNorm.Y*rightDist},
		}
	}

	// check top wall
	topDist := dot(tNorm.X, c.X-tLine.From.X, tNorm.Y, c.Y-tLine.From.Y)
	if topDist < c.Radius && movingToward(c, tNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: tNorm,
			Depth:  c.Radius - topDist,
			Point:  Vector{X: c.X - tNorm.X*topDist, Y: c.Y - tNorm.Y*topDist},
		}
	}

	// check bottom wall
	bottomDist := dot(bNorm.X, c.X-bLine.From.X, bNorm.Y, c.Y-bLine.From.Y)
	if bottomDist < c.Radius && movingToward(c, bNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: bNorm,
			Depth:  c.Radius - bottomDist,
			Point:  Vector{X: c.X - bNorm.X*bottomDist, Y: c.Y - bNorm.Y*bottomDist},
		}
	}

	return Collision{}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
