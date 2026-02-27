package levels

import "math"

type Point struct {
	X float32
	Y float32
}

func (p Point) Add(v Point) Point {
	return Point{X: p.X + v.X, Y: p.Y + v.Y}
}

func (p Point) Sub(v Point) Point {
	return Point{X: p.X - v.X, Y: p.Y - v.Y}
}

func (p Point) Dot(u Vector) float32 {
	return p.X*u.X + p.Y*u.Y
}

func (p Point) Scale(scalar float32) Point {
	return Point{X: p.X * scalar, Y: p.Y * scalar}
}

func (p Point) RotateAround(origin Point, angle float32) Point {
	sin, cos := math.Sincos(float64(angle))
	translated := p.Sub(origin)
	rotated := Point{
		X: translated.X*float32(cos) - translated.Y*float32(sin),
		Y: translated.X*float32(sin) + translated.Y*float32(cos),
	}
	return rotated.Add(origin)
}

type Size struct {
	W float32
	H float32
}

func (s Size) Scale(scalar float32) Size {
	return Size{W: s.W * scalar, H: s.H * scalar}
}

type Vector Point

func (v Vector) Add(u Vector) Vector {
	return Vector{X: v.X + u.X, Y: v.Y + u.Y}
}

func (v Vector) Reflect(normal Vector) Vector {
	dotProduct := v.Dot(normal)
	return Vector{
		X: v.X - 2*dotProduct*normal.X,
		Y: v.Y - 2*dotProduct*normal.Y,
	}
}

func (v Vector) Sub(u Vector) Vector {
	return Vector{X: v.X - u.X, Y: v.Y - u.Y}
}

func (v Vector) Scale(scalar float32) Vector {
	return Vector{X: v.X * scalar, Y: v.Y * scalar}
}

func (v Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vector) Dot(u Vector) float32 {
	return v.X*u.X + v.Y*u.Y
}

func (v Vector) Normalize() Vector {
	length := v.Length()
	if length == 0 {
		return Vector{X: 0, Y: 0} // arbitrary fallback if zero length
	}
	return v.Scale(1 / length)
}
