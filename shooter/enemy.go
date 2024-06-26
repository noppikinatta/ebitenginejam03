package shooter

import (
	"math"
	"math/rand/v2"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type EnemyLauncher struct {
	Enemies     []*Enemy
	Speed       float64
	FirstWait   int
	Rnd         *rand.Rand
	StageSize   geom.PointF
	ShipHit     geom.Circle
	Interval    int
	CurrentWait int
	Annihilated bool
}

func (l *EnemyLauncher) Update() {
	for _, e := range l.Enemies {
		e.Update()
	}

	if l.CurrentWait > 0 {
		l.CurrentWait--
		return
	}
	l.CurrentWait = int(float64(l.Interval) * (l.Rnd.Float64()*0.4 + 0.8))
	l.launch()
}

func (l *EnemyLauncher) launch() {
	enemyRemaining := false
	for _, e := range l.Enemies {
		if e.State != EnemyStateReady {
			if e.State == EnemyStateOnStage {
				enemyRemaining = true
			}
			continue
		}

		start := l.startPos()
		e.Launch(start, l.Speed, l.ShipHit, l.FirstWait)
		return // annihilated flag is not set
	}

	if !enemyRemaining {
		l.Annihilated = true
	}
}

func (l *EnemyLauncher) startPos() geom.PointF {
	ratio := l.Rnd.Float64()
	var zeroOrMax float64 = 0
	if l.Rnd.Float64() < 0.5 {
		zeroOrMax = 1
	}

	if l.Rnd.Float64() < 0.5 {
		return geom.PointF{
			X: ratio * l.StageSize.X,
			Y: zeroOrMax * l.StageSize.Y,
		}
	} else {
		return geom.PointF{
			X: zeroOrMax * l.StageSize.X,
			Y: ratio * l.StageSize.Y,
		}
	}
}

func (l *EnemyLauncher) Bullets() []Bullet {
	bb := make([]Bullet, 0)
	for _, e := range l.Enemies {
		for _, b := range e.Bullets {
			bb = append(bb, b)
		}
	}
	return bb
}

func (l *EnemyLauncher) Targets() []Target {
	tt := make([]Target, 0)
	for _, e := range l.Enemies {
		tt = append(tt, e)
	}
	return tt
}

type Enemy struct {
	HP               int
	MaxHP            int
	State            EnemyState
	Hit              geom.Circle
	Speed            float64
	Clockwise        bool
	ShootingInterval int
	CurrentWait      int
	ShipHit          geom.Circle
	MinCloseToShip   float64
	Bullets          []*EnemyBullet
	Rnd              *rand.Rand
}

func (e *Enemy) Launch(start geom.PointF, speed float64, shipHit geom.Circle, firstWait int) {
	e.State = EnemyStateOnStage
	e.Hit.Center = start
	e.Speed = speed
	e.CurrentWait = firstWait
	e.ShipHit = shipHit
	e.Clockwise = e.Rnd.Float64() < 0.5
}

func (e *Enemy) Update() {
	for _, b := range e.Bullets {
		b.Update()
	}

	if e.State != EnemyStateOnStage {
		return
	}

	if e.HP <= 0 {
		e.State = EnemyStateDead
		return
	}

	e.Hit.Center = e.Hit.Center.Add(e.velocity())
	e.shoot()
}

func (e *Enemy) velocity() geom.PointF {
	velocity := geom.PointF{}

	vecToShip := e.ShipHit.Center.Subtract(e.Hit.Center)
	allowedDistance := vecToShip.Abs() - (e.MinCloseToShip + e.ShipHit.Radius)
	if allowedDistance > 0 {
		velocity = geom.PointFFromPolar(e.Speed*(allowedDistance/e.ShipHit.Radius), vecToShip.Angle())
	}

	if e.Clockwise {
		velocity = velocity.Add(geom.PointFFromPolar(e.Speed, vecToShip.Angle()-0.5*math.Pi))
	} else {
		velocity = velocity.Add(geom.PointFFromPolar(e.Speed, vecToShip.Angle()+0.5*math.Pi))
	}

	return velocity
}

func (e *Enemy) shoot() {
	if e.CurrentWait > 0 {
		e.CurrentWait--
		return
	}

	shot := false
	for _, b := range e.Bullets {
		if b.IsLiving() {
			continue
		}

		b.Shoot(e.bulletInitParams())
		shot = true
		break
	}

	interval := e.ShootingInterval
	if shot {
		interval /= 4
	}

	e.CurrentWait = interval
}

func (e *Enemy) bulletInitParams() (start, velocity geom.PointF) {
	angle := e.ShipHit.Center.Subtract(e.Hit.Center).Angle()
	angle += (e.Rnd.Float64()*10 - 5) * math.Pi / 180
	abs := e.Speed * 4

	return e.Hit.Center, geom.PointFFromPolar(abs, angle)
}

func (e *Enemy) HitCircle() geom.Circle {
	return e.Hit
}

func (e *Enemy) IsEnemy() bool {
	return true
}

func (e *Enemy) Damage(value int) float64 {
	e.HP -= value
	if e.HP <= 0 {
		e.State = EnemyStateDead
		return e.Hit.Radius * 1.5
	}
	return 0
}

func (e *Enemy) IsLiving() bool {
	return e.State == EnemyStateOnStage
}

type EnemyState int

const (
	EnemyStateReady EnemyState = iota
	EnemyStateOnStage
	EnemyStateDead
)

type EnemyBullet struct {
	Power    int
	Hit      geom.Circle
	Velocity geom.PointF
	Cruising bool
}

func (b *EnemyBullet) Shoot(start, velocity geom.PointF) {
	b.Cruising = true
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
	return b.Cruising
}

func (b *EnemyBullet) HitProcess(targets []Target) geom.Circle {
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

		r := t.Damage(b.Power)
		b.Cruising = false

		if r > 0 {
			return geom.Circle{Center: b.Hit.Center, Radius: float64(b.Power)}
		} else {
			break
		}
	}

	return geom.Circle{}
}
