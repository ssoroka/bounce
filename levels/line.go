package levels

import "math"

type Line struct {
	From Point
	To   Point
}

func (l Line) Length() float32 {
	dx := l.To.X - l.From.X
	dy := l.To.Y - l.From.Y
	return (float32)(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (l Line) Normalized() Vector {
	length := l.Length()
	if length == 0 {
		return Vector{X: 0, Y: 0}
	}
	return Vector{X: (l.To.X - l.From.X) / length, Y: (l.To.Y - l.From.Y) / length}
}

func (l Line) Dot(v Vector) float32 {
	lineVec := Vector{X: l.To.X - l.From.X, Y: l.To.Y - l.From.Y}
	return lineVec.Dot(v)
}

func (l Line) Normal() Vector {
	normalVec := Vector{X: l.To.Y - l.From.Y, Y: l.From.X - l.To.X} // Rotate 90 degrees
	return normalVec.Normalize()
}

func (l Line) Intersect(other Line) (Point, bool) {
	dx1 := l.To.X - l.From.X
	dy1 := l.To.Y - l.From.Y
	dx2 := other.To.X - other.From.X
	dy2 := other.To.Y - other.From.Y

	denom := dx1*dy2 - dy1*dx2
	if denom == 0 {
		return Point{}, false // Parallel
	}

	dx3 := other.From.X - l.From.X
	dy3 := other.From.Y - l.From.Y

	t := (dx3*dy2 - dy3*dx2) / denom
	u := (dx3*dy1 - dy3*dx1) / denom

	// t and u must both be in [0, 1] for intersection to be on both segments
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return Point{}, false
	}

	return Point{
		X: l.From.X + t*dx1,
		Y: l.From.Y + t*dy1,
	}, true
}

// ContainsPoint checks if a point lies on the line segment
func (l Line) ContainsPoint(p Point) bool {
	const eps = 0.0001

	// Check if point is within bounding box
	minX := min(l.From.X, l.To.X) - eps
	maxX := max(l.From.X, l.To.X) + eps
	minY := min(l.From.Y, l.To.Y) - eps
	maxY := max(l.From.Y, l.To.Y) + eps

	if p.X < minX || p.X > maxX || p.Y < minY || p.Y > maxY {
		return false
	}

	// Check if point is collinear using cross product
	// (p - From) × (To - From) should be ~0
	dx := l.To.X - l.From.X
	dy := l.To.Y - l.From.Y
	px := p.X - l.From.X
	py := p.Y - l.From.Y

	cross := dx*py - dy*px
	return cross > -eps && cross < eps
}
