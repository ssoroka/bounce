package levels

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fps       = 80
	itemCount = 50
)

type Game struct {
	Window   Size
	Objects  []Drawable
	LastTick time.Time
}

type Drawable interface {
	Draw(s *ebiten.Image)
	Update(lastTick time.Time) error
}

func (g *Game) Update() (err error) {
	for _, o := range g.Objects {
		if err = o.Update(g.LastTick); err != nil {
			return err
		}
	}

	g.LastTick = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, c := range g.Objects {
		c.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.Window.W), int(g.Window.H)
}

func Level1() {
	g := &Game{
		Window: Size{W: 1710, H: 1112},
	}
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(g.Window.W), int(g.Window.H))
	for range itemCount {
		size := rand.Float32() * 80
		g.Objects = append(g.Objects, &Circle{
			Point: Point{rand.Float32() * g.Window.W, rand.Float32() * g.Window.H},
			// Size:     Size{size, size},
			Radius:   size / 2,
			Color:    color.RGBA{0, 0, 255, 255},
			Velocity: Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1},
		})
	}
	g.Objects = append(g.Objects, &Boundary{
		Point:       Point{0, 0},
		Size:        Size{g.Window.W - 2, g.Window.H - 2},
		StrokeWidth: 2,
		Color:       color.RGBA{255, 0, 255, 255},
	})
	g.LastTick = time.Now()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
