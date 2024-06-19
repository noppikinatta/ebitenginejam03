package build

import "github.com/noppikinatta/ebitenginejam03/geom"

type DecisionMaker struct {
	Left     float64
	Length   float64
	LinearFn geom.LinearFunc
}

func NewDecisionMaker(left float64, length float64) *DecisionMaker {
	return &DecisionMaker{
		Left:     left,
		Length:   length,
		LinearFn: geom.LinearFunc{A: 0, B: 1, C: 0}, // Y = 0
	}
}

func (d *DecisionMaker) Update(xCenter float64) {
	halfLen := d.Length / 2
	newX1 := xCenter - halfLen
	d.Left = newX1
}
