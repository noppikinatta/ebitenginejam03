package shooter

import "github.com/noppikinatta/ebitenginejam03/geom"

type Enemy struct {
	HP       int
	State    State
	Hit      geom.Circle
	Velocity geom.PointF
	Rotate   float64
	Bullets  []*EnemyBullet
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
	// TODO: update bullets
}

type EnemyBullet struct {
	Hit      geom.Circle
	Velocity geom.PointF
	State    State
}
