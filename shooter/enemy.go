package shooter

import (
	"math"
	"math/rand/v2"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type Enemy struct {
	HP               int
	State            State
	Hit              geom.Circle
	Velocity         geom.PointF
	ShootingInterval int
	CurrentWait      int
	Bullets          []*EnemyBullet
	Rnd              *rand.Rand
}

func (e *Enemy) Update() {
	if e.State != StateOnStage {
		return
	}

	if e.HP <= 0 {
		e.State = StateDead
		return
	}

	e.Hit.Center = e.Hit.Center.Add(e.Velocity)
	e.shoot()
}

func (e *Enemy) shoot() {
	if e.CurrentWait > 0 {
		e.CurrentWait--
		return
	}

	shot := false
	for _, b := range e.Bullets {
		b.Update()
		if b.IsLiving() {
			continue
		}

		if !shot {
			b.Shoot(e.bulletInitParams())
			shot = true
		}
	}

	interval := e.ShootingInterval
	if shot {
		interval /= 4
	}

	e.CurrentWait = interval
}

func (e *Enemy) bulletInitParams() (start, velocity geom.PointF) {
	angle := e.Velocity.Angle()
	angle += (e.Rnd.Float64()*10 - 5) * math.Pi / 180
	abs := e.Velocity.Abs() * 4

	return e.Hit.Center, geom.PointFFromPolar(abs, angle)
}

func (e *Enemy) HitCircle() geom.Circle {
	return e.Hit
}

func (e *Enemy) IsEnemy() bool {
	return true
}

func (e *Enemy) Damage(value int) {
	e.HP -= value
	if e.HP <= 0 {
		e.State = StateDead
	}
}

func (e *Enemy) IsLiving() bool {
	return e.State == StateOnStage
}

type EnemyBullet struct {
	Power    int
	Hit      geom.Circle
	Velocity geom.PointF
	State    State
}

func (b *EnemyBullet) Shoot(start, velocity geom.PointF) {
	b.State = StateOnStage
	b.Hit.Center = start
	b.Velocity = velocity
}

func (b *EnemyBullet) Update() {
	if !b.IsLiving() {
		return
	}
	b.Hit.Center = b.Hit.Center.Add(b.Velocity)
}

func (b *EnemyBullet) IsLiving() bool {
	return b.State == StateOnStage
}

func (b *EnemyBullet) HitProcess(targets []Target) {
	for _, t := range targets {
		if !t.IsLiving() {
			continue
		}
		if t.IsEnemy() {
			continue
		}

		if !b.Hit.IntersectsWith(t.HitCircle()) {
			continue
		}

		t.Damage(b.Power)
		b.State = StateReady
		break
	}
}
