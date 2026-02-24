package levels

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Boundary struct {
	Point
	Size
	StrokeWidth float32
	Color       color.Color
}

func (b *Boundary) Draw(screen *ebiten.Image) {
	vector.StrokeRect(screen, b.X, b.Y, b.W, b.H, b.StrokeWidth, b.Color, true)
}

func (b *Boundary) Update(lastTick time.Time) error {
	return nil
}
