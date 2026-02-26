package levels

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Cube struct {
	Point
	Size
	Color    color.Color
	Velocity Vector
}

func (c *Cube) Draw(s *ebiten.Image) {
	vector.FillRect(s, c.X, c.Y, c.W, c.H, c.Color, true)
}

func (c *Cube) Update(delta float32) error {
	c.X = c.X + c.Velocity.X*delta*fps
	c.Y = c.Y + c.Velocity.Y*delta*fps

	return nil
}
