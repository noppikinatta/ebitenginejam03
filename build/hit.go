package build

type HitBox struct {
	Position       PointF
	Size           PointF
	Velocity       PointF
	Rotate         float64
	RotateVelocity float64
}

func (b *HitBox) Left() float64 {
	return b.Position.X
}

func (b *HitBox) Right() float64 {
	return b.Position.X + b.Size.X
}

func (b *HitBox) Top() float64 {
	return b.Position.Y
}

func (b *HitBox) Bottom() float64 {
	return b.Position.Y + b.Size.Y
}

func (b *HitBox) Center() PointF {
	return PointF{
		X: b.Position.X + b.Size.X*0.5,
		Y: b.Position.Y + b.Size.Y*0.5,
	}
}

func (b *HitBox) Update() {
	b.Position.X += b.Velocity.X
	b.Position.Y += b.Velocity.Y
	b.Rotate += b.RotateVelocity
}

func (b *HitBox) BoundTop(y float64) {
	if b.Position.Y < y {
		b.Position.Y = y - b.Position.Y
	}
	if b.Velocity.Y < 0 {
		b.Velocity.Y *= -1
	}
}

func (b *HitBox) BoundBottom(y float64) {
	if b.Position.Y > y {
		b.Position.Y = 2*y - b.Position.Y
	}
	if b.Velocity.Y > 0 {
		b.Velocity.Y *= -1
	}
}

func (b *HitBox) BoundLeft(x float64) {
	if b.Position.X < x {
		b.Position.X = 2*x - b.Position.X
	}
	if b.Velocity.X < 0 {
		b.Velocity.X *= -1
	}
}

func (b *HitBox) BoundRight(x float64) {
	if b.Position.X > x {
		b.Position.X = 2*x - b.Position.X
	}
	if b.Velocity.X > 0 {
		b.Velocity.X *= -1
	}
}

func (b *HitBox) MultiplyVelocity(v float64) {
	a := b.Velocity.Abs()
	r := b.Velocity.Direction360()

	a *= v

	b.Velocity = PointFFromPolar(a, r)
}

func (b *HitBox) AddRotateVelocity(v float64) {
	b.RotateVelocity += v
}

type HorizontalHitLine struct {
	Position PointF
	Length   float64
}

func (l *HorizontalHitLine) Left() float64 {
	return l.Position.X
}

func (l *HorizontalHitLine) Right() float64 {
	return l.Position.X + l.Length
}

func (l *HorizontalHitLine) InXRange(x float64) bool {
	return x >= l.Left() && x <= l.Right()
}

func (l *HorizontalHitLine) JustCrossed(pt1, pt2 PointF) bool {
	// above of line
	if pt1.Y < l.Position.Y && pt2.Y < l.Position.Y {
		return false
	}

	// below of line
	if pt1.Y >= l.Position.Y && pt2.Y >= l.Position.Y {
		return false
	}

	// special case: vertical line
	if pt1.X == pt2.X {
		return l.InXRange(pt1.X)
	}

	// linear line except vertical
	f := LinearFuncFromPt(pt1, pt2)
	x := f.X(l.Position.Y)
	return l.InXRange(x)
}

func (l *HorizontalHitLine) IntersectsWith(box HitBox) bool {
	if box.Left() < l.Left() {
		return false
	}
	if box.Right() > l.Right() {
		return false
	}
	if box.Top() > l.Position.Y {
		return false
	}
	if box.Bottom() < l.Position.Y {
		return false
	}

	return true
}
