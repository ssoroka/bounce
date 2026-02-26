package levels

import "math"

func dot(x1, x2, y1, y2 float32) float32 {
	return x1*x2 + y1*y2
}

func normal(p1, p2 Point) (Line, Vector) {
	midPoint := Point{X: (p2.X + p1.X) / 2.0, Y: (p2.Y + p1.Y) / 2.0}

	dy := p2.Y - p1.Y
	dx := p2.X - p1.X
	dx, dy = dy, -dx // rotate 90 degrees counterclockwise to get normal vector

	// normalize the normal vector
	length := (float32)(math.Sqrt(float64(dx*dx + dy*dy)))
	mag := float32(1.0) / length
	dx *= mag
	dy *= mag

	normalLine := Line{
		From: Point{X: midPoint.X, Y: midPoint.Y},
		To:   Point{X: dx*50 + midPoint.X, Y: dy*50 + midPoint.Y},
	}

	return normalLine, Vector{X: dx, Y: dy}
}

// func crossProd
