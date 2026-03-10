package levels

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	fps       = 60
	itemCount = 5
	FPSDelta  = 1 / float32(fps)
)

type GameOptions struct {
	Fullscreen bool
}

type Game struct {
	Window         Size
	WindowPosition Point
	Objects        []Drawable
	LastTick       time.Time
	Options        GameOptions
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

	g.CheckKeyboardInput()
	g.ApplyGravity()
	g.CheckCollisions()

	g.LastTick = g.LastTick.Add(deltaDur)
	return nil
}

var (
	debug          = false
	gravity        = false
	collisionCount = 0
)

func (g *Game) Draw(screen *ebiten.Image) {
	for _, c := range g.Objects {
		c.Draw(screen)
	}

	if drawing {
		switch currentDrawObject {
		case DrawObjectBoundary:
			vector.StrokeLine(screen, drawStart.X, drawStart.Y, drawEnd.X, drawEnd.Y, 2, purple, false)
		case DrawObjectCube:
			size := Point{X: drawEnd.X - drawStart.X, Y: drawEnd.Y - drawStart.Y}
			pos := Point{X: drawStart.X, Y: drawStart.Y}
			c := NewCube(pos.X, pos.Y, size.X, size.Y, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}, Vector{0, 0})
			c.Draw(screen)
		case DrawObjectCircle:
			radius := Vector{X: drawEnd.X - drawStart.X, Y: drawEnd.Y - drawStart.Y}.Length()
			c := NewCircle(drawStart.X, drawStart.Y, radius, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}, Vector{0, 0})
			c.Draw(screen)
		}
	}

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(`FPS: %.2f
Velocity: %.2f
Count: %d
Collisions: %d
Velocity Init: %t
Current Draw Object %s`, ebiten.ActualFPS(), velocity, len(g.Objects)-1, collisionCount, initWithVelocity, currentDrawObject.String()))
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

func (g *Game) CheckCollisions() {
	for i, o1 := range g.Objects {
		for j, o2 := range g.Objects {
			if i >= j {
				continue
			}
			if col := CheckCollision(o1, o2); col.Hit {
				collisionCount++
				// d := col.Depth
				// d := float32(1.0)
				if _, ok := o2.(*CubeBoundary); ok {
					o1, o2 = o2, o1
					col.Depth = -col.Depth
				}
				if _, ok := o2.(*Boundary); ok {
					o1, o2 = o2, o1
					col.Depth = -col.Depth
				}
				// if _, ok := o1.(*Cube); ok {
				// 	if _, ok := o2.(*Circle); ok {
				// 		o1, o2 = o2, o1
				// 		col.Depth = -col.Depth
				// 	}
				// }
				// if one is a boundary, push the other out
				if _, ok := o1.(*CubeBoundary); ok {
					if c, ok := o2.(*Circle); ok {
						// push o2
						c.Velocity = c.Velocity.Reflect(col.Normal)

						// Push circle out of wall (adjust position, not velocity)
						c.X += col.Normal.X * col.Depth
						c.Y += col.Normal.Y * col.Depth
					}
					// if c, ok := o2.(*Cube); ok {
					// 	// push o2
					// 	c.Velocity = c.Velocity.Reflect(col.Normal)

					// 	// Push cube out of wall (adjust position, not velocity)
					// 	c.X += col.Normal.X * col.Depth
					// 	c.Y += col.Normal.Y * col.Depth
					// }
				} else if _, ok := o1.(*Boundary); ok {
					if c, ok := o2.(*Circle); ok {
						fmt.Println("wall collision!")
						// push o2
						c.Velocity = c.Velocity.Reflect(col.Normal)

						// Push circle out of wall (adjust position, not velocity)
						c.X += col.Normal.X * col.Depth
						c.Y += col.Normal.Y * col.Depth
					}
				} else if c1, ok := o1.(*Circle); ok {
					if c2, ok := o2.(*Circle); ok {
						// two circles
						// c1 := o1.(*Circle)
						// c2 := o2.(*Circle)

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
					}
					// if c2, ok := o2.(*Cube); ok {
					// 	// circle and cube
					// 	_ = c1
					// 	_ = c2
					// }
					// } else if c1, ok := o1.(*Cube); ok {
					// 	if c2, ok := o2.(*Cube); ok {
					// 		// two cubes
					// 		_ = c1
					// 		_ = c2
					// 	}
				}
			}
		}
	}
}

type DrawObjectType int

const (
	DrawObjectBoundary DrawObjectType = iota
	DrawObjectCube
	DrawObjectCircle
)

func (t DrawObjectType) String() string {
	switch t {
	case DrawObjectBoundary:
		return "Boundary"
	case DrawObjectCube:
		return "Cube"
	case DrawObjectCircle:
		return "Circle"
	default:
		return "Unknown"
	}
}

var (
	drawing           = false
	drawStart         Point
	drawEnd           Point
	currentDrawObject = DrawObjectBoundary
	initWithVelocity  = true
)

func (g *Game) CheckKeyboardInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		gravity = !gravity
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
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
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || (ebiten.IsKeyPressed(ebiten.KeyZ) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		for i, obj := range g.Objects {
			switch obj.(type) {
			case *Boundary, *CubeBoundary:
				continue
			}
			// delete this item
			g.Objects = append(g.Objects[:i], g.Objects[i+1:]...)
			break
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || (ebiten.IsKeyPressed(ebiten.KeyX) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		createCube(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) || (ebiten.IsKeyPressed(ebiten.KeyC) && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		createCircle(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		debug = !debug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if !recording {
			if err := g.StartRecording(int(g.Window.W), int(g.Window.H)); err != nil {
				log.Println("error starting recording:", err)
			}
		} else {
			filename := g.StopRecording()
			log.Println("recording saved to", filename)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if err := g.SaveState("save.json"); err != nil {
			log.Println("error saving state:", err)
		} else {
			log.Println("state saved to save.json")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		if err := g.LoadState("save.json"); err != nil {
			log.Println("error loading state:", err)
		} else {
			log.Println("state loaded from save.json")
		}
	}

	// drawing
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		currentDrawObject = DrawObjectBoundary
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		currentDrawObject = DrawObjectCube
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		currentDrawObject = DrawObjectCircle
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		initWithVelocity = !initWithVelocity
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.Options.Fullscreen = !g.Options.Fullscreen
		ebiten.SetFullscreen(g.Options.Fullscreen)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if drawing {
			return
		}
		drawing = true
		x, y := ebiten.CursorPosition()
		drawStart = Point{X: float32(x), Y: float32(y)}
	}
	if drawing && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		drawing = false
		x, y := ebiten.CursorPosition()
		drawEnd = Point{X: float32(x), Y: float32(y)}
		switch currentDrawObject {
		case DrawObjectBoundary:
			g.Objects = append(g.Objects, NewBoundaryLine(drawStart, drawEnd, 2, purple))
		case DrawObjectCube:
			size := Point{X: drawEnd.X - drawStart.X, Y: drawEnd.Y - drawStart.Y}
			pos := Point{X: drawStart.X, Y: drawStart.Y}
			var velocity Vector
			if initWithVelocity {
				velocity = Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1}
			}
			c := NewCube(pos.X, pos.Y, size.X, size.Y, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}, velocity)
			g.Objects = append(g.Objects, c)
		case DrawObjectCircle:
			radius := Vector{X: drawEnd.X - drawStart.X, Y: drawEnd.Y - drawStart.Y}.Length()
			var velocity Vector
			if initWithVelocity {
				velocity = Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1}
			}
			c := NewCircle(drawStart.X, drawStart.Y, radius, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}, velocity)
			g.Objects = append(g.Objects, c)
		}
	} else if drawing {
		x, y := ebiten.CursorPosition()
		drawEnd = Point{X: float32(x), Y: float32(y)}
	}
}

var (
	cube   *CubeBoundary
	pixels []byte
)

func Level1() {
	windowW, windowH := ebiten.Monitor().Size()
	g := &Game{
		Options: GameOptions{Fullscreen: true},
		Window:  Size{W: float32(windowW), H: float32(windowH)},
	}

	// make a buffer for reading pixels when recording
	pixels = make([]byte, int(g.Window.W)*int(g.Window.H)*4)

	ebiten.SetTPS(fps)
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetFullscreen(g.Options.Fullscreen)
	ebiten.SetWindowSize(int(g.Window.W), int(g.Window.H))
	cube = NewCubeBoundary(0, 0, g.Window.W-2, g.Window.H-2, 2, purple)
	// g.Objects = append(g.Objects, cube)
	boundary := NewBoundary(0, 0, g.Window.W-2, g.Window.H-2, 2, purple)
	g.Objects = append(g.Objects, boundary)
	g.Objects = append(g.Objects, createCube(g))
	for range itemCount {
		createCircle(g)
	}

	g.LastTick = time.Now()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func createCircle(g *Game) {
	size := rand.Float32()*20 + 10
	// make sure it's a random point that fits within the bounds of the rotated cube
	x := rand.Float32() * (cube.W - size)
	y := rand.Float32() * (cube.H - size)
	color := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}
	velocity := Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1}
	g.Objects = append(g.Objects, NewCircle(x, y, size/2, color, velocity))
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

func createCube(g *Game) *Cube {
	size := rand.Float32()*70 + 10
	x := rand.Float32() * (cube.W - size)
	y := rand.Float32() * (cube.H - size)
	pos := Point{x, y}.RotateAround(Point{X: cube.W / 2, Y: cube.H / 2}, cube.Rotation).Add(Point{X: cube.X, Y: cube.Y})
	c := NewCube(pos.X, pos.Y, size, size, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}, Vector{rand.Float32()*2 - 1, rand.Float32()*2 - 1})
	g.Objects = append(g.Objects, c)
	return c
}

var gravityConstant = float32(9.8 / fps * 2)

func (g *Game) ApplyGravity() {
	if !gravity {
		return
	}
	for _, o := range g.Objects {
		if c, ok := o.(*Circle); ok {
			c.Velocity.Y += gravityConstant
		}
	}
}
