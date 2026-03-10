package levels

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Spring struct {
	c1 *Circle
	c2 *Circle

	Length    float32
	Stiffness float32
	Thickness float32
	Color     color.Color
}

func NewSpring(from, to *Circle, stiffness, thickness float32, color color.Color) *Spring {
	dx := to.X - from.X
	dy := to.Y - from.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	return &Spring{
		c1:        from,
		c2:        to,
		Length:    distance,
		Stiffness: stiffness,
		Thickness: thickness,
		Color:     color,
	}
}

func (s *Spring) Update(delta float32) error {
	dx := s.c2.X - s.c1.X
	dy := s.c2.Y - s.c1.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance == 0 {
		return nil
	}

	// Calculate the force magnitude based on Hooke's law
	forceMagnitude := s.Stiffness * (distance - s.Length)

	// Normalize the direction vector
	nx := dx / distance
	ny := dy / distance

	// Apply the force to each circle (equal and opposite)
	s.c1.Velocity.X += forceMagnitude * nx
	s.c1.Velocity.Y += forceMagnitude * ny
	s.c2.Velocity.X -= forceMagnitude * nx
	s.c2.Velocity.Y -= forceMagnitude * ny
	return nil
}

func (s *Spring) Draw(surf *ebiten.Image) {
	vector.StrokeLine(surf, s.c1.X, s.c1.Y, s.c2.X, s.c2.Y, s.Thickness, s.Color, true)
}
