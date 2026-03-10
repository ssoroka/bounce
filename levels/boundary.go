package levels

import (
	"encoding/json"
	"image/color"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CubeBoundary struct {
	Point
	Size
	Rotation    float32
	StrokeWidth float32
	Color       color.Color
	tl          Point
	tr          Point
	bl          Point
	br          Point
}

func NewCubeBoundary(x, y, w, h, strokeWidth float32, color color.Color) *CubeBoundary {
	c := &CubeBoundary{
		Point:       Point{X: x, Y: y},
		Size:        Size{W: w, H: h},
		StrokeWidth: strokeWidth,
		Color:       color,
	}
	c.RecalculateCorners()
	return c
}

func (b *CubeBoundary) Scale(factor float32) {
	// Scale the size, but keep the center of the cube the same
	centerX := b.X + b.W/2
	centerY := b.Y + b.H/2
	b.X = centerX - (b.W*factor)/2
	b.Y = centerY - (b.H*factor)/2
	b.Size = b.Size.Scale(factor)
	b.RecalculateCorners()
}

func (b *CubeBoundary) GetEdges() [4]Line {
	return [4]Line{
		{From: b.tl, To: b.bl},
		{From: b.bl, To: b.br},
		{From: b.br, To: b.tr},
		{From: b.tr, To: b.tl},
	}
}

func (b *CubeBoundary) RecalculateCorners() {
	// Center of rotation
	cx := b.X + b.W/2
	cy := b.Y + b.H/2

	// Unrotated corners relative to center
	corners := []Point{
		{-b.W / 2, -b.H / 2}, // tl
		{b.W / 2, -b.H / 2},  // tr
		{-b.W / 2, b.H / 2},  // bl
		{b.W / 2, b.H / 2},   // br
	}

	sin, cos := float32(math.Sin(float64(b.Rotation))), float32(math.Cos(float64(b.Rotation)))

	rotate := func(p Point) Point {
		return Point{
			X: cx + p.X*cos - p.Y*sin,
			Y: cy + p.X*sin + p.Y*cos,
		}
	}

	b.tl = rotate(corners[0])
	b.tr = rotate(corners[1])
	b.bl = rotate(corners[2])
	b.br = rotate(corners[3])
}

func (b *CubeBoundary) Draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, b.tl.X, b.tl.Y, b.tr.X, b.tr.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.tl.X, b.tl.Y, b.bl.X, b.bl.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.bl.X, b.bl.Y, b.br.X, b.br.Y, b.StrokeWidth, b.Color, true)
	vector.StrokeLine(screen, b.tr.X, b.tr.Y, b.br.X, b.br.Y, b.StrokeWidth, b.Color, true)

	// draw normals
	if debug {
		lLine, lNorm := normal(b.tl, b.bl)
		bLine, bNorm := normal(b.bl, b.br)
		rLine, rNorm := normal(b.br, b.tr)
		tLine, tNorm := normal(b.tr, b.tl)

		vector.StrokeLine(screen, lLine.From.X, lLine.From.Y, lLine.To.X, lLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, bLine.From.X, bLine.From.Y, bLine.To.X, bLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, rLine.From.X, rLine.From.Y, rLine.To.X, rLine.To.Y, b.StrokeWidth, green, true)
		vector.StrokeLine(screen, tLine.From.X, tLine.From.Y, tLine.To.X, tLine.To.Y, b.StrokeWidth, green, true)
		// dot(lNorm.X, 0, lNorm.Y, 0)
		_ = lNorm
		_ = bNorm
		_ = rNorm
		_ = tNorm

	}
}

func (b *CubeBoundary) Update(delta float32) error {
	return nil
}

func movingToward(c *Circle, norm Vector) bool {
	return dot(c.Velocity.X, norm.X, c.Velocity.Y, norm.Y) < 0
}

// CheckCircleCollision checks if a circle trapped within a cube boundary is colliding with the
// boundary walls and returns collision info to keep it inside.
func (b *CubeBoundary) CheckCircleCollision(c *Circle) Collision {
	lLine, lNorm := normal(b.tl, b.bl)
	bLine, bNorm := normal(b.bl, b.br)
	rLine, rNorm := normal(b.br, b.tr)
	tLine, tNorm := normal(b.tr, b.tl)

	// check left wall
	leftDist := dot(lNorm.X, c.X-lLine.From.X, lNorm.Y, c.Y-lLine.From.Y)
	if leftDist < c.Radius && movingToward(c, lNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: lNorm,
			Depth:  c.Radius - leftDist,
			Point:  Vector{X: c.X - lNorm.X*leftDist, Y: c.Y - lNorm.Y*leftDist},
		}
	}

	// check right wall
	rightDist := dot(rNorm.X, c.X-rLine.From.X, rNorm.Y, c.Y-rLine.From.Y)
	if rightDist < c.Radius && movingToward(c, rNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: rNorm,
			Depth:  c.Radius - rightDist,
			Point:  Vector{X: c.X - rNorm.X*rightDist, Y: c.Y - rNorm.Y*rightDist},
		}
	}

	// check top wall
	topDist := dot(tNorm.X, c.X-tLine.From.X, tNorm.Y, c.Y-tLine.From.Y)
	if topDist < c.Radius && movingToward(c, tNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: tNorm,
			Depth:  c.Radius - topDist,
			Point:  Vector{X: c.X - tNorm.X*topDist, Y: c.Y - tNorm.Y*topDist},
		}
	}

	// check bottom wall
	bottomDist := dot(bNorm.X, c.X-bLine.From.X, bNorm.Y, c.Y-bLine.From.Y)
	if bottomDist < c.Radius && movingToward(c, bNorm) {
		return Collision{
			Hit:    true, // only consider it a collision if the circle is moving towards the wall
			Normal: bNorm,
			Depth:  c.Radius - bottomDist,
			Point:  Vector{X: c.X - bNorm.X*bottomDist, Y: c.Y - bNorm.Y*bottomDist},
		}
	}

	return Collision{}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

type Boundary struct {
	Lines       []Line
	StrokeWidth float32
	Color       color.Color
}

const (
	pointCount = 6
	variance   = float32(100)
)

func NewBoundaryLine(from, to Point, strokeWidth float32, color color.Color) *Boundary {
	return &Boundary{
		Lines:       []Line{{From: from, To: to}},
		StrokeWidth: strokeWidth,
		Color:       color,
	}
}

func (b *Boundary) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type        string  `json:"type"`
		Lines       []Line  `json:"lines"`
		StrokeWidth float32 `json:"strokeWidth"`
		ColorR      uint8   `json:"R"`
		ColorG      uint8   `json:"G"`
		ColorB      uint8   `json:"B"`
		ColorA      uint8   `json:"A"`
	}{
		Type:        "Boundary",
		Lines:       b.Lines,
		StrokeWidth: b.StrokeWidth,
		ColorR:      uint8(b.Color.(color.RGBA).R),
		ColorG:      uint8(b.Color.(color.RGBA).G),
		ColorB:      uint8(b.Color.(color.RGBA).B),
		ColorA:      uint8(b.Color.(color.RGBA).A),
	})
}

func (b *Boundary) UnmarshalJSON(data []byte) error {
	aux := struct {
		Type        string  `json:"type"`
		Lines       []Line  `json:"lines"`
		StrokeWidth float32 `json:"strokeWidth"`
		ColorR      uint8   `json:"R"`
		ColorG      uint8   `json:"G"`
		ColorB      uint8   `json:"B"`
		ColorA      uint8   `json:"A"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	b.Lines = aux.Lines
	b.StrokeWidth = aux.StrokeWidth
	b.Color = color.RGBA{R: aux.ColorR, G: aux.ColorG, B: aux.ColorB, A: aux.ColorA}
	return nil
}

func NewBoundary(x, y, w, h, strokeWidth float32, color color.Color) *Boundary {
	lines := []Line{}

	// LEFT WALL
	// pick 3 points inside the volume area to connect to, somewhat close to the wall
	start := Point{x, y}
	end := Point{x, y + h}
	juts := make([]Point, pointCount)
	for i := range pointCount {
		juts[i] = Point{
			X: x + rand.Float32()*variance,
			Y: rand.Float32()*(end.Y-start.Y) + start.Y,
		}
	}
	sort.Slice(juts, func(i, j int) bool {
		return juts[i].Y < juts[j].Y
	})

	// create lines between the points and the start and end
	appendLines := func(juts []Point, start, end Point) {
		for i := range pointCount {
			if i == 0 {
				lines = append(lines, Line{From: start, To: juts[i]})
			} else if i == pointCount-1 {
				lines = append(lines, Line{From: juts[i-1], To: juts[i]})
				lines = append(lines, Line{From: juts[i], To: end})
			} else {
				lines = append(lines, Line{From: juts[i-1], To: juts[i]})
			}
		}
	}

	appendLines(juts, start, end)

	// BOTTOM WALL
	// pick 3 points inside the volume area to connect to, somewhat close to the wall
	start = Point{x, y + h}
	end = Point{x + w, y + h}
	for i := range pointCount {
		juts[i] = Point{X: x + rand.Float32()*(end.X-start.X+end.X), Y: start.Y - rand.Float32()*variance}
	}
	sort.Slice(juts, func(i, j int) bool {
		return juts[i].X < juts[j].X
	})

	// create lines between the points and the start and end
	appendLines(juts, start, end)

	// RIGHT WALL
	// pick 3 points inside the volume area to connect to, somewhat close to the wall
	start = Point{x + w, y + h}
	end = Point{x + w, y}
	for i := range pointCount {
		juts[i] = Point{start.X - rand.Float32()*variance, rand.Float32()*(end.Y-start.Y) + start.Y}
	}
	sort.Slice(juts, func(i, j int) bool {
		return juts[i].Y > juts[j].Y
	})

	// create lines between the points and the start and end
	appendLines(juts, start, end)

	// TOP WALL
	start = Point{x + w, y}
	end = Point{x, y}
	for i := range pointCount {
		juts[i] = Point{rand.Float32()*(end.X-start.X) + start.X, start.Y + rand.Float32()*variance}
	}
	sort.Slice(juts, func(i, j int) bool {
		return juts[i].X > juts[j].X
	})

	// create lines between the points and the start and end
	appendLines(juts, start, end)

	return &Boundary{Lines: lines, StrokeWidth: strokeWidth, Color: color}
}

func (b *Boundary) Draw(screen *ebiten.Image) {
	for _, line := range b.Lines {
		vector.StrokeLine(screen, line.From.X, line.From.Y, line.To.X, line.To.Y, 2, purple, true)
		lLine, _ := normal(line.From, line.To)
		if debug {
			vector.StrokeLine(screen, lLine.From.X, lLine.From.Y, lLine.To.X, lLine.To.Y, b.StrokeWidth, green, true)
		}
	}
}

func (b *Boundary) Update(delta float32) error {
	return nil
}

// func (b *Boundary) CheckCircleCollision(c *Circle) Collision {
// 	for _, line := range b.Lines {
// 		// check collision with each line segment
// 		lLine, lNorm := normal(line.From, line.To)
// 		rotatedNormal := Vector{X: -lNorm.Y, Y: lNorm.X}

// 		toCircle := Vector{X: c.X - line.From.X, Y: c.Y - line.From.Y}

// 		rotatedDistToLine := dot(rotatedNormal.X, c.X-line.From.X, rotatedNormal.Y, c.Y-line.From.Y)

// 		lineDist := dot(lNorm.X, c.X-lLine.From.X, lNorm.Y, c.Y-lLine.From.Y)

// 		withinLine := rotatedDistToLine >= 0 && rotatedDistToLine <= line.Length() // ✅ Within [0, length]
// 		// withinLine := math.Abs(float64(rotatedDistToLine)) < float64(line.Length())/2

// 		if lineDist > 0 && lineDist < c.Radius && withinLine && movingToward(c, lNorm) {
// 			return Collision{
// 				Hit:    true, // only consider it a collision if the circle is moving towards the wall
// 				Normal: lNorm,
// 				Depth:  c.Radius - lineDist,
// 				Point:  Vector{X: c.X - lNorm.X*lineDist, Y: c.Y - lNorm.Y*lineDist},
// 			}
// 		}
// 	}
// 	return Collision{}
// }

func (b *Boundary) CheckCircleCollision(c *Circle) Collision {
	for _, line := range b.Lines {
		// Find closest point on line segment to circle center
		closestPoint := line.ClosestPoint(Vector{X: c.X, Y: c.Y})

		// Distance from circle center to closest point
		delta := Vector{X: c.X - closestPoint.X, Y: c.Y - closestPoint.Y}
		dist := delta.Length()

		normal := delta.Normalize()
		if dist < c.Radius && c.Velocity.Dot(normal.Scale(-1)) > 0 {
			return Collision{
				Hit:    true,
				Normal: normal,
				Depth:  c.Radius - dist,
				Point:  closestPoint,
			}
		}
	}

	startPos := c.LastPosition
	endPos := Point{X: c.X, Y: c.Y}

	for _, line := range b.Lines {
		// Check if circle path intersects line
		if t := raySegmentIntersect(startPos, endPos, c.Radius, line); t >= 0 {
			// Collision at time t along movement
			collisionPos := startPos.Add(endPos.Sub(startPos).Scale(t))
			// calculate collision response: bounce off the line normal
			lineVec := Vector{X: line.To.X - line.From.X, Y: line.To.Y - line.From.Y}
			normal := Vector{X: -lineVec.Y, Y: lineVec.X}.Normalize()
			if c.Velocity.Dot(normal.Scale(-1)) > 0 {
				return Collision{
					Hit:    true,
					Normal: normal,
					Depth:  c.Radius, // could be refined to actual penetration depth
					Point:  Vector(collisionPos),
				}
			}
		}
	}

	return Collision{}
}

func raySegmentIntersect(start, end Point, radius float32, line Line) float32 {
	// Ray from start to end
	rayDir := Vector{X: end.X - start.X, Y: end.Y - start.Y}
	rayLength := rayDir.Length()
	if rayLength == 0 {
		return -1 // No movement
	}
	rayDir = rayDir.Scale(1 / rayLength) // Normalize

	// Line segment vector
	lineDir := Vector{X: line.To.X - line.From.X, Y: line.To.Y - line.From.Y}
	lineLength := lineDir.Length()
	if lineLength == 0 {
		return -1 // Invalid line
	}
	lineDir = lineDir.Scale(1 / lineLength) // Normalize

	// Calculate intersection using cross products
	cross := rayDir.X*lineDir.Y - rayDir.Y*lineDir.X
	if math.Abs(float64(cross)) < 1e-8 {
		return -1 // Parallel lines
	}

	diff := Vector{X: line.From.X - start.X, Y: line.From.Y - start.Y}
	t := (diff.X*lineDir.Y - diff.Y*lineDir.X) / cross
	u := (diff.X*rayDir.Y - diff.Y*rayDir.X) / cross

	if t >= 0 && t <= rayLength && u >= 0 && u <= lineLength {
		return t / rayLength // Return normalized time of collision
	}

	return -1 // No collision
}

func (l Line) ClosestPoint(p Vector) Vector {
	lineVec := Vector{X: l.To.X - l.From.X, Y: l.To.Y - l.From.Y}
	toPoint := Vector{X: p.X - l.From.X, Y: p.Y - l.From.Y}

	t := toPoint.Dot(lineVec) / lineVec.Dot(lineVec)
	t = clamp(t, 0, 1) // Clamp to segment bounds

	return Vector{
		X: l.From.X + t*lineVec.X,
		Y: l.From.Y + t*lineVec.Y,
	}
}
