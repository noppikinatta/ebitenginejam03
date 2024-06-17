package build

import "math"

type PointF struct {
	X, Y float64
}

func (p PointF) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p PointF) DirectionRad() float64 {
	return math.Atan2(p.Y, p.X)
}

func (p PointF) Direction360() float64 {
	return 360 * p.DirectionRad() / math.Pi
}

func PointFFromPolar(a float64, r360 float64) PointF {
	rRad := r360 * math.Pi / 180.0

	x := a * math.Cos(rRad)
	y := a * math.Sin(rRad)

	return PointF{X: x, Y: y}
}

type LinearFunc struct {
	Slope, Intercept float64
}

func LinearFuncFromPt(pt1, pt2 PointF) LinearFunc {
	slope := (pt2.Y - pt1.Y) / (pt2.X - pt1.X)
	intercept := pt1.Y - slope*pt1.X

	return LinearFunc{Slope: slope, Intercept: intercept}
}

func (f LinearFunc) Y(x float64) float64 {
	return x*f.Slope + f.Intercept
}

func (f LinearFunc) X(y float64) float64 {
	return (y - f.Intercept) / f.Slope
}
