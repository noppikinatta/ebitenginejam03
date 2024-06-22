package shooter

import (
	"math"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type Equip struct {
	Name     string
	Position geom.PointF
	Angle    float64
	Updater  EquipUpdater
}

func (e *Equip) UpdateAngle(shipHit geom.Circle, angle float64) {
	pos := geom.PointFFromPolar(shipHit.Radius, angle)
	pos = pos.Add(shipHit.Center)
	e.Position = pos
	e.Angle = angle
}

func (e *Equip) Update(enemies []*Enemy) {
	e.Updater.Update(e, enemies)
}

type EquipUpdater interface {
	Update(equip *Equip, enemies []*Enemy)
}

type EquipUpdaterNop struct{}

func (u *EquipUpdaterNop) Update(equip *Equip, enemies []*Enemy) {}

type EquipUpdaterLaser struct {
	ShipHit     geom.Circle
	LastFrames  int
	Interval    int
	CurrentWait int
	CurrentLast int
	Width       float64
	Power       int
}

func (u *EquipUpdaterLaser) Update(equip *Equip, enemies []*Enemy) {
	if u.CurrentWait > 0 {
		u.CurrentWait--
		if u.CurrentWait == 0 && u.CurrentLast == 0 {
			u.CurrentLast = u.LastFrames
		}
		return
	}

	if u.CurrentLast == 0 {
		u.CurrentWait = u.Interval
		return
	}

	u.CurrentLast--

	line := geom.LinearFuncFromPt(u.ShipHit.Center, equip.Position)
	for _, enemy := range enemies {
		distance := line.Distance(enemy.Hit.Center)
		if distance > u.Width {
			return
		}
		enemy.HP -= u.Power
	}
}

func (u *EquipUpdaterLaser) LaserLastingRatio() float64 {
	return float64(u.CurrentLast) / float64(u.LastFrames)
}

type EquipUpdaterMissile struct {
	Interval    int
	CurrentWait int
	Missiles    []*Missile
}

func (u *EquipUpdaterMissile) Update(equip *Equip, enemies []*Enemy) {
	for _, m := range u.Missiles {
		if m.State != StateOnStage {
			continue
		}

		m.Update(enemies)
	}

	if u.CurrentWait > 0 {
		u.CurrentWait--
		return
	}

	u.CurrentWait = u.Interval
	for _, m := range u.Missiles {
		if m.State == StateOnStage {
			continue
		}

		u.launchMissile(equip, m)
		break
	}
}

func (u *EquipUpdaterMissile) launchMissile(equip *Equip, missile *Missile) {
	missile.Hit.Center = equip.Position
	missile.Velocity = geom.PointFFromPolar(missile.Velocity.Abs(), equip.Angle)
	missile.State = StateOnStage
}

type Missile struct {
	Hit         geom.Circle
	Velocity    geom.PointF
	PowerToTurn float64
	State       State
	Power       int
}

func (m *Missile) Update(enemies []*Enemy) {
	var closestEnemy *Enemy
	var closestDistance float64 = math.Inf(1)
	exploded := false
	for _, e := range enemies {
		if e.State != StateOnStage {
			continue
		}

		if m.Hit.IntersectsWith(e.Hit) {
			exploded = true
			e.HP -= m.Power
			continue
		}

		distance := m.Hit.Center.Distance(e.Hit.Center)
		if distance < closestDistance {
			closestDistance = distance
			closestEnemy = e
		}
	}

	if exploded {
		m.State = StateDead
		return
	}

	if closestEnemy == nil {
		// There may be no enemies at all
		return
	}

	vectorForEnemy := closestEnemy.Hit.Center.Add(m.Hit.Center.Multiply(-1))
	enemyAngle := math.Atan2(vectorForEnemy.X, vectorForEnemy.Y)
	currentAngle := m.Velocity.Angle()
	gapAngle := currentAngle - enemyAngle
	if gapAngle < -1*m.PowerToTurn {
		gapAngle = -1 * m.PowerToTurn
	}
	if gapAngle > m.PowerToTurn {
		gapAngle = m.PowerToTurn
	}
	m.Velocity = geom.PointFFromPolar(m.Velocity.Abs(), currentAngle-gapAngle)
}

type EquipUpdaterHarakiriSystem struct {
	MyShipHit   geom.Circle
	Interval    int
	CurrentWait int
	Harakiris   []*HarakiriSystem
}

func (u *EquipUpdaterHarakiriSystem) Update(equip *Equip, enemies []*Enemy) {
	for _, m := range u.Harakiris {
		if m.State != StateOnStage {
			continue
		}

		m.Update(enemies)
	}

	if u.CurrentWait > 0 {
		u.CurrentWait--
		return
	}

	u.CurrentWait = u.Interval
	for _, m := range u.Harakiris {
		if m.State == StateOnStage {
			continue
		}

		u.launchMissile(equip, m)
		break
	}

}

type HarakiriSystem struct {
	Hit      geom.Circle
	Velocity geom.PointF
	State    State
	Power    int
}

func (h *HarakiriSystem) Update(enemies []*Enemy, myShipHit geom.Circle) {
	// TODO: destroy when hit barriers
}

type EquipUpdaterBarrier struct {
	Radius      float64
	Count       int
	Max         int
	Interval    int
	CurrentWait int
}

func (u *EquipUpdaterBarrier) Update(equip *Equip, enemies []*Enemy) {

}

type EquipUpdaterExhaust struct {
	Radius     float64
	Multiplier float64
}

func (u *EquipUpdaterExhaust) Update(equip *Equip, enemies []*Enemy) {

}
