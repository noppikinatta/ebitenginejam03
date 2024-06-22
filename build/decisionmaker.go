package build

import (
	"github.com/noppikinatta/ebitenginejam03/geom"
)

type DecisionMaker struct {
	LeftX    float64
	Width    float64
	Y        float64
	LinearFn geom.LinearFunc
}

func NewDecisionMaker(leftX float64, width float64) *DecisionMaker {
	var y float64 = 0

	d := &DecisionMaker{
		LeftX: leftX,
		Width: width,
		Y:     y,
	}

	d.LinearFn = geom.LinearFuncFromPt(d.Left(), d.Right())
	return d
}

func (d *DecisionMaker) Left() geom.PointF {
	return geom.PointF{X: d.LeftX, Y: d.Y}
}

func (d *DecisionMaker) Right() geom.PointF {
	return geom.PointF{X: d.LeftX + d.Width, Y: d.Y}
}

func (d *DecisionMaker) Update(xCenter float64) {
	halfLen := d.Width / 2
	newX1 := xCenter - halfLen
	d.LeftX = newX1
}

func (d *DecisionMaker) Hit(c geom.Circle) bool {
	if c.Center.Y < d.Y {
		return false // proposal passed decision maker
	}

	dist := d.LinearFn.Distance(c.Center)
	if dist > c.Radius {
		return false
	}

	if c.Left() > d.Right().X {
		return false
	}

	if c.Right() < d.Left().X {
		return false
	}

	return true
}
