package shooter

import (
	"math"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type MyShip struct {
	HP     int
	Hit    geom.Circle
	Angle  float64
	Equips []*Equip
}

func (m *MyShip) UpdateAngle(angle float64) {
	m.Angle = angle

	l := len(m.Equips)
	if l == 0 {
		return
	}
	pi2 := math.Pi * 2

	for i, e := range m.Equips {
		angle := pi2*float64(i)/float64(l) + m.Angle
		e.UpdateAngle(m.Hit, angle)
	}
}

func (m *MyShip) Update(enemies []*Enemy) {
	for _, e := range m.Equips {
		e.Update(enemies)
	}
}
