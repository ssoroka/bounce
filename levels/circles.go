package levels

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Circle struct {
	Point
	Radius   float32
	Color    color.Color
	Velocity Vector
}

func (c *Circle) Draw(s *ebiten.Image) {
	vector.FillCircle(s, c.X, c.Y, c.Radius, c.Color, true)
}

func (c *Circle) Update(lastTick time.Time) error {
	delta := float32(time.Since(lastTick).Seconds())

	c.X = c.X + c.Velocity.X*delta*fps
	c.Y = c.Y + c.Velocity.Y*delta*fps
	return nil
}
