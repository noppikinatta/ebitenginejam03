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

func (m *MyShip) Update() {
	for _, e := range m.Equips {
		e.Update()
	}
}

func (m *MyShip) HitCircle() geom.Circle {
	return m.Hit
}

func (m *MyShip) IsEnemy() bool {
	return false
}

func (m *MyShip) Damage(value int) float64 {
	m.HP -= value
	return 0
}

func (m *MyShip) IsLiving() bool {
	return m.HP > 0
}

func (m *MyShip) Bullets() []Bullet {
	bb := make([]Bullet, 0)
	for _, e := range m.Equips {
		bb = append(bb, e.Updater.Bullets()...)
	}
	return bb
}

func (m *MyShip) Targets() []Target {
	tt := make([]Target, 0)
	for _, e := range m.Equips {
		tt = append(tt, e.Updater.Targets()...)
	}
	tt = append(tt, m)
	return tt
}

func (m *MyShip) VisibleEntities() []VisibleEntity {
	vv := make([]VisibleEntity, 0)
	for _, e := range m.Equips {
		vv = append(vv, e.Updater.VisibleEntities()...)
	}
	return vv
}
