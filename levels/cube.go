package levels

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	X float32
	Y float32
}

type Size struct {
	W float32
	H float32
}

type Cube struct {
	Point
	Size
	Color    color.Color
	Velocity Vector
}

type Vector Point

func (c *Cube) Draw(s *ebiten.Image) {
	vector.FillRect(s, c.X, c.Y, c.W, c.H, c.Color, true)
}

func (c *Cube) Update(lastTick time.Time) error {
	delta := float32(time.Since(lastTick).Seconds())

	c.X = c.X + c.Velocity.X*delta*fps
	c.Y = c.Y + c.Velocity.Y*delta*fps

	return nil
}
