package levels

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Cube struct {
	Point
	Size
	Rotation float32
	Color    color.Color
	Velocity Vector
	tl       Point
	tr       Point
	bl       Point
	br       Point
}

func NewCube(x, y, w, h float32, color color.Color) *Cube {
	c := &Cube{
		Point: Point{X: x, Y: y},
		Size:  Size{W: w, H: h},
		Color: color,
	}
	c.RecalculateCorners()
	return c
}

func (c *Cube) Draw(s *ebiten.Image) {
	vector.FillRect(s, c.X, c.Y, c.W, c.H, c.Color, true)

	// draw normals
	if debug {
		tLine, lNorm := normal(c.tl, c.tr)
		rLine, bNorm := normal(c.tr, c.br)
		bLine, rNorm := normal(c.br, c.bl)
		lLine, tNorm := normal(c.bl, c.tl)

		vector.StrokeLine(s, lLine.From.X, lLine.From.Y, lLine.To.X, lLine.To.Y, 1, green, true)
		vector.StrokeLine(s, bLine.From.X, bLine.From.Y, bLine.To.X, bLine.To.Y, 1, green, true)
		vector.StrokeLine(s, rLine.From.X, rLine.From.Y, rLine.To.X, rLine.To.Y, 1, green, true)
		vector.StrokeLine(s, tLine.From.X, tLine.From.Y, tLine.To.X, tLine.To.Y, 1, green, true)
		// dot(lNorm.X, 0, lNorm.Y, 0)
		_ = lNorm
		_ = bNorm
		_ = rNorm
		_ = tNorm

	}
}

func (c *Cube) Update(delta float32) error {
	c.X = c.X + c.Velocity.X
	c.Y = c.Y + c.Velocity.Y
	c.RecalculateCorners()

	return nil
}

// approx radius to a corner used for distance checks before more expensive SAT collision checks
func (c *Cube) Radius() float32 {
	return c.W / 2 * float32(math.Sqrt2)
}

func (c *Cube) Center() Point {
	return Point{X: c.X + c.W/2, Y: c.Y + c.H/2}
}

func (c *Cube) RecalculateCorners() {
	// Center of rotation
	cx := c.X + c.W/2
	cy := c.Y + c.H/2

	// Unrotated corners relative to center
	corners := []Point{
		{-c.W / 2, -c.H / 2}, // tl
		{c.W / 2, -c.H / 2},  // tr
		{-c.W / 2, c.H / 2},  // bl
		{c.W / 2, c.H / 2},   // br
	}

	sin, cos := float32(math.Sin(float64(c.Rotation))), float32(math.Cos(float64(c.Rotation)))

	rotate := func(p Point) Point {
		return Point{
			X: cx + p.X*cos - p.Y*sin,
			Y: cy + p.X*sin + p.Y*cos,
		}
	}

	c.tl = rotate(corners[0])
	c.tr = rotate(corners[1])
	c.bl = rotate(corners[2])
	c.br = rotate(corners[3])
}

func (c *Cube) GetCorners() [4]Vector {
	// // Same rotation logic as your CubeBoundary
	// cx, cy := c.X+c.W/2, c.Y+c.H/2
	// hw, hh := c.W/2, c.H/2
	// sin, cos := float32(math.Sin(float64(c.Rotation))), float32(math.Cos(float64(c.Rotation)))

	// return [4]Vector{
	// 	{cx + (-hw)*cos - (-hh)*sin, cy + (-hw)*sin + (-hh)*cos}, // tl
	// 	{cx + (hw)*cos - (-hh)*sin, cy + (hw)*sin + (-hh)*cos},   // tr
	// 	{cx + (hw)*cos - (hh)*sin, cy + (hw)*sin + (hh)*cos},     // br
	// 	{cx + (-hw)*cos - (hh)*sin, cy + (-hw)*sin + (hh)*cos},   // bl
	// }
	return [4]Vector{
		{X: c.tl.X, Y: c.tl.Y},
		{X: c.tr.X, Y: c.tr.Y},
		{X: c.br.X, Y: c.br.Y},
		{X: c.bl.X, Y: c.bl.Y},
	}
}

func (c *Cube) GetLines() [4]Line {
	return [4]Line{
		{From: c.tl, To: c.tr},
		{From: c.tr, To: c.br},
		{From: c.br, To: c.bl},
		{From: c.bl, To: c.tl},
	}
}

func (c *Cube) GetAxes() [2]Vector {
	// Edge normals (perpendicular to edges)
	sin, cos := float32(math.Sin(float64(c.Rotation))), float32(math.Cos(float64(c.Rotation)))
	return [2]Vector{
		{cos, sin},  // perpendicular to top/bottom edges
		{-sin, cos}, // perpendicular to left/right edges
	}
}
