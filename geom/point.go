package geom

import "math"

type PointF struct {
	X, Y float64
}

func (p PointF) Add(other PointF) PointF {
	return PointF{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p PointF) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p PointF) DirectionRad() float64 {
	return math.Atan2(p.Y, p.X)
}

func (p PointF) Distance(other PointF) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func PointFFromPolar(a float64, rRad float64) PointF {
	x := a * math.Cos(rRad)
	y := a * math.Sin(rRad)

	return PointF{X: x, Y: y}
}
