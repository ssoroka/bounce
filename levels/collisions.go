package levels

type Collision struct {
	Hit    bool
	Normal Vector  // direction to push objects apart
	Depth  float32 // penetration depth
	Point  Vector  // contact point (optional, useful for effects)
}

func CheckCollision(a, b any) Collision {
	switch sa := a.(type) {
	case *Circle:
		switch sb := b.(type) {
		case *Circle:
			return CircleVsCircle(sa, sb)
		case *CubeBoundary:
			return CircleVsRect(sa, sb)
		}
	case *CubeBoundary:
		switch sb := b.(type) {
		case *Circle:
			c := CircleVsRect(sb, sa)
			c.Normal = c.Normal.Scale(-1) // flip normal
			return c
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

func CircleVsRect(c *Circle, r *CubeBoundary) Collision {
	return r.CheckCircleCollision(c)
}
