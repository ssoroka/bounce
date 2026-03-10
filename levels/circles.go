package levels

import (
	"encoding/json"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Circle struct {
	Point
	LastPosition Point
	Radius       float32
	Color        color.Color
	Velocity     Vector
}

func NewCircle(x, y, radius float32, color color.Color, velocity Vector) *Circle {
	c := &Circle{
		Point:        Point{X: x, Y: y},
		Radius:       radius,
		Color:        color,
		Velocity:     velocity,
		LastPosition: Point{X: x, Y: y},
	}
	return c
}

func (c *Circle) Draw(s *ebiten.Image) {
	vector.FillCircle(s, c.X, c.Y, c.Radius, c.Color, true)
}

func (c *Circle) Update(delta float32) error {
	c.LastPosition = c.Point
	c.X = c.X + c.Velocity.X
	c.Y = c.Y + c.Velocity.Y
	return nil
}

func (c *Circle) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string  `json:"type"`
		Point    Point   `json:"point"`
		Radius   float32 `json:"radius"`
		ColorR   uint8   `json:"R"`
		ColorG   uint8   `json:"G"`
		ColorB   uint8   `json:"B"`
		ColorA   uint8   `json:"A"`
		Velocity Vector  `json:"velocity"`
	}{
		Type:     "Circle",
		Point:    c.Point,
		Radius:   c.Radius,
		ColorR:   uint8(c.Color.(color.RGBA).R),
		ColorG:   uint8(c.Color.(color.RGBA).G),
		ColorB:   uint8(c.Color.(color.RGBA).B),
		ColorA:   uint8(c.Color.(color.RGBA).A),
		Velocity: c.Velocity,
	})
}

func (c *Circle) UnmarshalJSON(data []byte) error {
	aux := struct {
		Type     string  `json:"type"`
		Point    Point   `json:"point"`
		Radius   float32 `json:"radius"`
		ColorR   uint8   `json:"R"`
		ColorG   uint8   `json:"G"`
		ColorB   uint8   `json:"B"`
		ColorA   uint8   `json:"A"`
		Velocity Vector  `json:"velocity"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Point = aux.Point
	c.Radius = aux.Radius
	c.Color = color.RGBA{R: aux.ColorR, G: aux.ColorG, B: aux.ColorB, A: aux.ColorA}
	c.Velocity = aux.Velocity

	return nil
}
