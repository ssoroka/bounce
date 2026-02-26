package levels

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	fps       = 100
	itemCount = 5
	FPSDelta  = 1 / float32(fps)
)

type Game struct {
	Window   Size
	Objects  []Drawable
	LastTick time.Time
}

type Drawable interface {
	Draw(s *ebiten.Image)
	Update(delta float32) error
}

var velocity = float32(0.0)

func (g *Game) Update() (err error) {
	deltaDur := time.Since(g.LastTick)
	// delta := float32(deltaDur.Seconds())
	delta := FPSDelta
	velocity = float32(0.0)

	for _, o := range g.Objects {
		if err = o.Update(delta); err != nil {
			return err
		}
		if o, ok := o.(*Circle); ok {
			velocity += o.Velocity.Length()
		}
	}

	g.CheckKeyboardInput(delta)
	g.CheckCollisions(delta)

	g.LastTick = g.LastTick.Add(deltaDur)
	return nil
}

var (
	debug          = false
	collisionCount = 0
)

func (g *Game) Draw(screen *ebiten.Image) {
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(`FPS: %.2f
Velocity: %.2f
Count: %d
Collisions: %d`, ebiten.ActualFPS(), velocity, len(g.Objects)-1, collisionCount))
	}

	for _, c := range g.Objects {
		c.Draw(screen)
	}

	if recording && ffmpegPipe != nil {
		screen.ReadPixels(pixels)
		if _, err := ffmpegPipe.Write(pixels); err != nil {
			log.Println("error writing to ffmpeg pipe:", err)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.Window.W), int(g.Window.H)
}

func (g *Game) CheckCollisions(delta float32) {
	for i, o1 := range g.Objects {
		for j, o2 := range g.Objects {
			if i >= j {
				continue
			}
			if col := CheckCollision(o1, o2); col.Hit {
				collisionCount++
				// d := col.Depth
				// d := float32(1.0)
				// if one is a boundary, push the other out
				if _, ok := o1.(*CubeBoundary); ok {
					// push o2
					c := o2.(*Circle)
					c.Velocity = c.Velocity.Reflect(col.Normal)

					// Push circle out of wall (adjust position, not velocity)
					c.X += col.Normal.X * col.Depth
					c.Y += col.Normal.Y * col.Depth

				} else if _, ok := o2.(*CubeBoundary); ok {
					// push o1
					c := o1.(*Circle)
					c.Velocity = c.Velocity.Reflect(col.Normal)

					// Push circle out of wall (adjust position, not velocity)
					c.X -= col.Normal.X * col.Depth
					c.Y -= col.Normal.Y * col.Depth
				} else {
					// two circles
					c1 := o1.(*Circle)
					c2 := o2.(*Circle)

					// Separate circles
					c1.X -= col.Normal.X * col.Depth / 2
					c1.Y -= col.Normal.Y * col.Depth / 2
					c2.X += col.Normal.X * col.Depth / 2
					c2.Y += col.Normal.Y * col.Depth / 2

					// Elastic collision (equal mass)
					// Swap velocity components along collision normal
					relVel := Vector{c1.Velocity.X - c2.Velocity.X, c1.Velocity.Y - c2.Velocity.Y}
					dot := relVel.X*col.Normal.X + relVel.Y*col.Normal.Y

					c1.Velocity.X -= dot * col.Normal.X
					c1.Velocity.Y -= dot * col.Normal.Y
					c2.Velocity.X += dot * col.Normal.X
					c2.Velocity.Y += dot * col.Normal.Y

					// c1.Velocity = c1.Velocity.Add(col.Normal.Scale(-col.Depth / 2.0))
					// c2.Velocity = c2.Velocity.Add(col.Normal.Scale(col.Depth / 2.0))

					// normalX := col.Normal.X
					// normalY := col.Normal.Y
					// relativeVelocityX := c1.Velocity.X - c2.Velocity.X
					// relativeVelocityY := c1.Velocity.Y - c2.Velocity.Y
					// relativeVelocityAlongNormal := dot(normalX, relativeVelocityX, normalY, relativeVelocityY)
					// relativeVelocityAlongTangent := dot(-normalY, relativeVelocityX, normalX, relativeVelocityY)

				}
			}
		}
	}
}

func (g *Game) CheckKeyboardInput(delta float32) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		cube.Scale(0.999)
		// cube.Size = cube.Size.Scale(0.99)
		// cube.RecalculateCorners()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		cube.Scale(1.001)
		// cube.Size = cube.Size.Scale(1.01)
		// cube.RecalculateCorners()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		cube.Rotation -= 0.001
		cube.RecalculateCorners()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		cube.Rotation += 0.001
		cube.RecalculateCorners()
	}
	if isKeyJustPressed(ebiten.KeyZ) || (ebiten.IsKeyPressed(ebiten.KeyZ) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		for i := range g.Objects {
			if g.Objects[i] == cube {
				continue
			}
			// delete this item
			g.Objects = append(g.Objects[:i], g.Objects[i+1:]...)
			break
		}
	}
	if isKeyJustPressed(ebiten.KeyX) || (ebiten.IsKeyPressed(ebiten.KeyX) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		createCircle(g)
	}
	if isKeyJustPressed(ebiten.KeyD) {
		debug = !debug
	}
	if isKeyJustPressed(ebiten.KeyR) {
		if !recording {
			if err := g.StartRecording(int(g.Window.W), int(g.Window.H)); err != nil {
				log.Println("error starting recording:", err)
			}
		} else {
			filename := g.StopRecording()
			log.Println("recording saved to", filename)
		}
	}
}

var (
	cube   *CubeBoundary
	pixels []byte
)

func Level1() {
	windowW, windowH := ebiten.Monitor().Size()
	g := &Game{
		Window: Size{W: float32(windowW), H: float32(windowH)},
	}

	// make a buffer for reading pixels when recording
	pixels = make([]byte, int(g.Window.W)*int(g.Window.H)*4)

	ebiten.SetTPS(fps)
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(g.Window.W), int(g.Window.H))
	cube = NewCubeBoundary(0, 0, g.Window.W-2, g.Window.H-2, 2, purple)
	g.Objects = append(g.Objects, cube)
	for range itemCount {
		createCircle(g)
	}

	g.LastTick = time.Now()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func createCircle(g *Game) {
	size := rand.Float32()*70 + 10
	// make sure it's a random point that fits within the bounds of the rotated cube
	x := rand.Float32() * (cube.W - size)
	y := rand.Float32() * (cube.H - size)
	c := &Circle{
		Point:    Point{x, y}.RotateAround(Point{X: cube.W / 2, Y: cube.H / 2}, cube.Rotation).Add(Point{X: cube.X, Y: cube.Y}),
		Radius:   size / 2,
		Color:    color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255},
		Velocity: Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1},
	}
	g.Objects = append(g.Objects, c)
}

var keyStates = make(map[ebiten.Key]bool)

func isKeyJustPressed(key ebiten.Key) bool {
	isPressed := ebiten.IsKeyPressed(key)
	if isPressed && !keyStates[key] {
		keyStates[key] = true
		return true
	} else if !isPressed {
		keyStates[key] = false
	}
	return false
}
