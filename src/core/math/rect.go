package math

import ()

// ============================================================================

type Rect struct {
	P1 Vector2
	P2 Vector2
}

// ============================================================================

func RectFromRange(pos Vector2, r int32) Rect {
	rv := Vector2{r, r}
	return Rect{pos.Sub(rv), pos.Add(rv)}
}

func (self *Rect) PtInside(pt Vector2) bool {
	return pt.X >= self.P1.X && pt.X <= self.P2.X && pt.Y >= self.P1.Y && pt.Y <= self.P2.Y
}

func (self *Rect) Bound(br Rect) {
	self.P1.X = MaxInt32(self.P1.X, br.P1.X)
	self.P1.Y = MaxInt32(self.P1.Y, br.P1.Y)

	self.P2.X = MinInt32(self.P2.X, br.P2.X)
	self.P2.Y = MinInt32(self.P2.Y, br.P2.Y)
}
