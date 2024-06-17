package build

type PointF struct {
	X, Y float64
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
