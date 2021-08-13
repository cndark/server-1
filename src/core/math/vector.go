package math

import (
	"math"
)

// ============================================================================

type Vector2 struct {
	X int32
	Y int32
}

// ============================================================================

func (a *Vector2) Set(b Vector2) {
	a.X = b.X
	a.Y = b.Y
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{a.X + b.X, a.Y + b.Y}
}

func (a Vector2) Sub(b Vector2) Vector2 {
	return Vector2{a.X - b.X, a.Y - b.Y}
}

func (a Vector2) Mul(v float32) Vector2 {
	return Vector2{int32(float32(a.X) * v), int32(float32(a.Y) * v)}
}

func (a Vector2) Dot(b Vector2) int32 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2) Cross(b Vector2) int32 {
	return a.X*b.Y - a.Y*b.X
}

func (a Vector2) Dist(b Vector2) float32 {
	return float32(math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))))
}

// rotate 90 degree
func (a Vector2) Rot90(ccw bool) {
	var tmp int32

	if ccw {
		tmp = -a.Y
		a.Y = a.X
		a.X = tmp
	} else {
		tmp = a.Y
		a.Y = -a.X
		a.X = tmp
	}
}

func (a Vector2) Equals(b Vector2) bool {
	return a.X == b.X && a.Y == b.Y
}
