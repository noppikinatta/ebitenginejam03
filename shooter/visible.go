package shooter

import "github.com/noppikinatta/ebitenginejam03/geom"

type VisibleEntity interface {
	Position() geom.PointF
	Angle() float64
	VisibleF() float64
	Name() string
}
