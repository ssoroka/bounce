package levels

import (
	"fmt"
	"math"
)

type Collision struct {
	Hit    bool
	Normal Vector  // direction to push objects apart
	Depth  float32 // penetration depth
	Point  Vector  // contact point (optional, useful for effects)
}

func CheckCollision(a, b any) Collision {
	switch sa := a.(type) {
	case *CubeBoundary:
		switch sb := b.(type) {
		case *Circle:
			c := CircleVsBoundary(sb, sa)
			c.Normal = c.Normal.Scale(-1) // flip normal
			return c
		case *Cube:
			return CubeVsBoundary(sb, sa)
		}
	case *Circle:
		switch sb := b.(type) {
		case *Circle:
			return CircleVsCircle(sa, sb)
		// case *CubeBoundary:
		// 	return CircleVsBoundary(sa, sb)
		case *Cube:
			return CircleVsCube(sa, sb)
		}
	case *Cube:
		switch sb := b.(type) {
		case *Circle:
			return CircleVsCube(sb, sa)
		case *Cube:
			return CubeVsCube(sa, sb)
		case *CubeBoundary:
			return CubeVsBoundary(sa, sb)
		}
	}
	return Collision{}
}

func CircleVsCircle(a, b *Circle) Collision {
	delta := Vector{X: b.X - a.X, Y: b.Y - a.Y}
	dist := delta.Length()
	minDist := a.Radius + b.Radius

	if dist >= minDist {
		return Collision{}
	}

	normal := delta.Normalize()

	return Collision{
		Hit:    true,
		Normal: normal,
		Depth:  minDist - dist,
		Point: Vector{
			X: a.X + normal.X*a.Radius,
			Y: a.Y + normal.Y*a.Radius,
		},
	}
}

func CircleVsBoundary(c *Circle, r *CubeBoundary) Collision {
	return r.CheckCircleCollision(c)
}

func CircleVsCube(c *Circle, cube *Cube) Collision {
	return Collision{}
}

func CubeVsCube(a, b *Cube) Collision {
	// if the cubes are axis-aligned, we can do a simple AABB check.
	if floatMod(a.Rotation, 90) == floatMod(b.Rotation, 90) {
		if a.X < b.X+b.W && a.X+a.W > b.X && a.Y < b.Y+b.H && a.Y+a.H > b.Y {
			// Simple AABB collision response (push apart along shortest axis)
			dx := min(a.X+a.W-b.X, b.X+b.W-a.X)
			dy := min(a.Y+a.H-b.Y, b.Y+b.H-a.Y)
			if dx < dy {
				if a.X < b.X {
					return Collision{Hit: true, Normal: Vector{X: -1, Y: 0}, Depth: dx}
				} else {
					return Collision{Hit: true, Normal: Vector{X: 1, Y: 0}, Depth: dx}
				}
			} else {
				if a.Y < b.Y {
					return Collision{Hit: true, Normal: Vector{X: 0, Y: -1}, Depth: dy}
				} else {
					return Collision{Hit: true, Normal: Vector{X: 0, Y: 1}, Depth: dy}
				}
			}
		}
		return Collision{}
	}
	// for rotated cubes, we use the Separating Axis Theorem (SAT)
	return CubeVsCubeSAT(a, b)

	// return Collision{}
}

func CubeVsBoundary(cube *Cube, boundary *CubeBoundary) Collision {
	// Use same approach as CheckCircleCollision but check cube corners
	edges := boundary.GetEdges()

	// For each boundary edge
	for edgeIdx, edge := range edges {
		edgeNormal := edge.Normal()

		// Check each corner of the cube
		for cornerIdx, corner := range cube.GetCorners() {
			// Distance from edge to this corner
			dist := dot(edgeNormal.X, corner.X-edge.From.X, edgeNormal.Y, corner.Y-edge.From.Y)

			// Debug: print for first cube
			if cornerIdx == 0 && edgeIdx == 0 {
				fmt.Printf("Edge %d, Corner %d: dist=%.2f, normal=(%.2f,%.2f), vel=(%.2f,%.2f)\n",
					edgeIdx, cornerIdx, dist, edgeNormal.X, edgeNormal.Y, cube.Velocity.X, cube.Velocity.Y)
			}

			// If corner is close to or past the wall (threshold of 1 pixel)
			if dist < 1 {
				// Check if cube is moving toward this wall
				velDot := cube.Velocity.X*edgeNormal.X + cube.Velocity.Y*edgeNormal.Y
				if velDot < 0 {
					depth := 1 - dist // how far to push back
					fmt.Printf("COLLISION! Edge %d, Corner %d: dist=%.2f, depth=%.2f\n",
						edgeIdx, cornerIdx, dist, depth)
					return Collision{
						Hit:    true,
						Normal: edgeNormal,
						Depth:  depth,
					}
				}
			}
		}
	}

	return Collision{}
}

func floatMod(a float32, b int) float32 {
	return a - float32(int(a)/b)*float32(b)
}

func CubeVsCubeSAT(a, b *Cube) Collision {
	// Get corners of both cubes
	cornersA := a.GetCorners() // [4]Vector (tl, tr, br, bl)
	cornersB := b.GetCorners()

	// Get edge normals (perpendicular to each edge)
	axesA := a.GetAxes() // [2]Vector (only need 2, opposite edges are parallel)
	axesB := b.GetAxes()

	axes := []Vector{axesA[0], axesA[1], axesB[0], axesB[1]}

	minDepth := float32(math.MaxFloat32)
	minAxis := Vector{}

	for _, axis := range axes {
		// Project all corners of A onto axis
		minA, maxA := projectCorners(cornersA, axis)
		// Project all corners of B onto axis
		minB, maxB := projectCorners(cornersB, axis)

		// Check for gap
		if maxA < minB || maxB < minA {
			return Collision{} // Separating axis found = no collision
		}

		// Calculate overlap depth on this axis
		overlap := min(maxA-minB, maxB-minA)
		if overlap < minDepth {
			minDepth = overlap
			minAxis = axis
		}
	}

	// No separating axis found = collision
	// minAxis is the collision normal, minDepth is penetration

	// Make sure normal points from A to B
	centerDiff := Vector{X: b.X - a.X, Y: b.Y - a.Y}
	if minAxis.Dot(centerDiff) < 0 {
		minAxis = minAxis.Scale(-1)
	}

	return Collision{
		Hit:    true,
		Normal: minAxis,
		Depth:  minDepth,
	}
}

func projectCorners(corners [4]Vector, axis Vector) (min, max float32) {
	min = corners[0].Dot(axis)
	max = min
	for i := 1; i < 4; i++ {
		p := corners[i].Dot(axis)
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}
	return
}
