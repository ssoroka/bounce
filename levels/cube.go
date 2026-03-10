package levels

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Cube struct {
	Points  []*Circle
	Springs []*Spring
	Size
	Color  color.Color
	Filled bool
}

func NewCube(x, y, w, h float32, color color.Color, velocity Vector) *Cube {
	// Unrotated corners relative to center
	corners := []Point{
		{x, y},         // tl
		{x, y + h},     // bl
		{x + w, y + h}, // br
		{x + w, y},     // tr
	}

	points := make([]*Circle, 4)
	springs := make([]*Spring, 4)

	for i, corner := range corners {
		points[i] = NewCircle(corner.X, corner.Y, 1, color, velocity)
	}

	springs[0] = NewSpring(points[0], points[1], 1, 1, color)
	springs[1] = NewSpring(points[1], points[2], 1, 1, color)
	springs[2] = NewSpring(points[2], points[3], 1, 1, color)
	springs[3] = NewSpring(points[3], points[0], 1, 1, color)

	return &Cube{
		Points:  points,
		Springs: springs,
		Size:    Size{W: w, H: h},
		Color:   color,
	}
}

func (c *Cube) Draw(s *ebiten.Image) {
	p := &vector.Path{}
	p.MoveTo(c.Points[0].X, c.Points[0].Y)
	p.LineTo(c.Points[1].X, c.Points[1].Y)
	p.LineTo(c.Points[2].X, c.Points[2].Y)
	p.LineTo(c.Points[3].X, c.Points[3].Y)
	p.Close()
	if c.Filled {
		vector.FillPath(s, p, &vector.FillOptions{}, &vector.DrawPathOptions{
			AntiAlias: true,
		})
	} else {
		vector.StrokePath(s, p, &vector.StrokeOptions{
			Width: 1,
		}, &vector.DrawPathOptions{
			AntiAlias: true,
		})
	}

	// draw normals
	if debug {
		tLine, _ := normal(c.Springs[0].c1.Point, c.Springs[0].c2.Point)
		rLine, _ := normal(c.Springs[1].c1.Point, c.Springs[1].c2.Point)
		bLine, _ := normal(c.Springs[2].c1.Point, c.Springs[2].c2.Point)
		lLine, _ := normal(c.Springs[3].c1.Point, c.Springs[3].c2.Point)

		vector.StrokeLine(s, lLine.From.X, lLine.From.Y, lLine.To.X, lLine.To.Y, 1, green, true)
		vector.StrokeLine(s, bLine.From.X, bLine.From.Y, bLine.To.X, bLine.To.Y, 1, green, true)
		vector.StrokeLine(s, rLine.From.X, rLine.From.Y, rLine.To.X, rLine.To.Y, 1, green, true)
		vector.StrokeLine(s, tLine.From.X, tLine.From.Y, tLine.To.X, tLine.To.Y, 1, green, true)
	}
}

func (c *Cube) Update(delta float32) error {
	for _, p := range c.Points {
		p.X = p.X + p.Velocity.X
		p.Y = p.Y + p.Velocity.Y
	}
	for _, s := range c.Springs {
		s.Update(delta)
	}
	return nil
}

// approx radius to a corner used for distance checks before more expensive SAT collision checks
// func (c *Cube) Radius() float32 {
// 	return c.W / 2 * float32(math.Sqrt2)
// }

// func (c *Cube) Center() Point {
// 	return Point{X: c.X + c.W/2, Y: c.Y + c.H/2}
// }

// func (c *Cube) RecalculateCorners() {
// 	// Center of rotation
// 	cx := c.X + c.W/2
// 	cy := c.Y + c.H/2

// 	// Unrotated corners relative to center
// 	corners := []Point{
// 		{-c.W / 2, -c.H / 2}, // tl
// 		{c.W / 2, -c.H / 2},  // tr
// 		{-c.W / 2, c.H / 2},  // bl
// 		{c.W / 2, c.H / 2},   // br
// 	}

// 	sin, cos := float32(math.Sin(float64(c.Rotation))), float32(math.Cos(float64(c.Rotation)))

// 	rotate := func(p Point) Point {
// 		return Point{
// 			X: cx + p.X*cos - p.Y*sin,
// 			Y: cy + p.X*sin + p.Y*cos,
// 		}
// 	}

// 	c.tl = rotate(corners[0])
// 	c.tr = rotate(corners[1])
// 	c.bl = rotate(corners[2])
// 	c.br = rotate(corners[3])
// }

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
		{X: c.Points[0].X, Y: c.Points[0].Y},
		{X: c.Points[1].X, Y: c.Points[1].Y},
		{X: c.Points[2].X, Y: c.Points[2].Y},
		{X: c.Points[3].X, Y: c.Points[3].Y},
	}
}

func (c *Cube) GetLines() [4]Line {
	return [4]Line{
		{From: c.Points[0].Point, To: c.Points[1].Point},
		{From: c.Points[1].Point, To: c.Points[2].Point},
		{From: c.Points[2].Point, To: c.Points[3].Point},
		{From: c.Points[3].Point, To: c.Points[0].Point},
	}
}

// func (c *Cube) GetAxes() [2]Vector {
// 	// Edge normals (perpendicular to edges)
// 	sin, cos := float32(math.Sin(float64(c.Rotation))), float32(math.Cos(float64(c.Rotation)))
// 	return [2]Vector{
// 		{cos, sin},  // perpendicular to top/bottom edges
// 		{-sin, cos}, // perpendicular to left/right edges
// 	}
// }
