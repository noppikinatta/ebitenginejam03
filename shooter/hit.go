package shooter

import "github.com/noppikinatta/ebitenginejam03/geom"

type Bullet interface {
	IsLiving() bool
	HitProcess(targets []Target)
}

type Target interface {
	HitCircle() geom.Circle
	IsEnemy() bool
	Damage(value int) float64
	IsLiving() bool
}

type HitTest struct {
	Bullets []Bullet
	Targets []Target
}

func (h *HitTest) Update() {
	for _, b := range h.Bullets {
		if !b.IsLiving() {
			continue
		}

		b.HitProcess(h.Targets)
	}
}
