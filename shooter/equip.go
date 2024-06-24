package shooter

import (
	"math"

	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
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

func (e *Equip) Update() {
	e.Updater.Update(e)
}

type EquipUpdater interface {
	Update(equip *Equip)
	Bullets() []Bullet
	Targets() []Target
	VisibleEntities() []VisibleEntity
}

type EquipUpdaterNop struct{}

func (u *EquipUpdaterNop) Update(equip *Equip)              {}
func (u *EquipUpdaterNop) Bullets() []Bullet                { return nil }
func (u *EquipUpdaterNop) Targets() []Target                { return nil }
func (u *EquipUpdaterNop) VisibleEntities() []VisibleEntity { return nil }

type EquipUpdaterLaser struct {
	ShipHit     geom.Circle
	Pos         geom.PointF
	LastFrames  int
	Interval    int
	CurrentWait int
	CurrentLast int
	Width       float64
	Power       int
}

func (u *EquipUpdaterLaser) Update(equip *Equip) {
	u.Pos = equip.Position

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
}

func (u *EquipUpdaterLaser) IsLiving() bool {
	return u.CurrentLast > 0
}

func (u *EquipUpdaterLaser) HitProcess(targets []Target) {
	line := geom.LinearFuncFromPt(u.ShipHit.Center, u.Pos)
	for _, target := range targets {
		if !target.IsEnemy() {
			continue
		}
		if !target.IsLiving() {
			continue
		}

		distance := line.Distance(target.HitCircle().Center)
		if distance > (u.Width + target.HitCircle().Radius) {
			continue
		}

		v1 := u.Pos.Subtract(u.ShipHit.Center)
		v2 := target.HitCircle().Center.Subtract(u.ShipHit.Center)
		if v1.InnerProduct(v2) < 0 {
			continue
		}

		target.Damage(u.Power)
	}
}

func (u *EquipUpdaterLaser) Position() geom.PointF {
	return u.Pos
}

func (u *EquipUpdaterLaser) Angle() float64 {
	return u.Pos.Subtract(u.ShipHit.Center).Angle()
}

func (u *EquipUpdaterLaser) VisibleF() float64 {
	return float64(u.CurrentLast) / float64(u.LastFrames)
}

func (u *EquipUpdaterLaser) Name() string {
	return name.EquipLaserCannon
}

func (u *EquipUpdaterLaser) Bullets() []Bullet {
	return []Bullet{u}
}

func (u *EquipUpdaterLaser) Targets() []Target {
	return nil
}

func (u *EquipUpdaterLaser) VisibleEntities() []VisibleEntity {
	return []VisibleEntity{u}
}

type EquipUpdaterMissile struct {
	Interval    int
	CurrentWait int
	Missiles    []*Missile
}

func (u *EquipUpdaterMissile) Update(equip *Equip) {
	for _, m := range u.Missiles {
		m.Update()
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
	missile.Launch(equip.Position, geom.PointFFromPolar(missile.FirstSpeed, equip.Angle))
}

func (u *EquipUpdaterMissile) Bullets() []Bullet {
	bb := make([]Bullet, len(u.Missiles))
	for i := range u.Missiles {
		bb[i] = u.Missiles[i]
	}
	return bb
}

func (u *EquipUpdaterMissile) Targets() []Target {
	tt := make([]Target, len(u.Missiles))
	for i := range u.Missiles {
		tt[i] = u.Missiles[i]
	}
	return tt
}

func (u *EquipUpdaterMissile) VisibleEntities() []VisibleEntity {
	vv := make([]VisibleEntity, len(u.Missiles))
	for i := range u.Missiles {
		vv[i] = u.Missiles[i]
	}
	return vv
}

type Missile struct {
	Hit          geom.Circle
	Velocity     geom.PointF
	FirstSpeed   float64
	Acceleration geom.PointF
	AccelPower   float64
	State        State
	Power        int
}

func (m *Missile) Update() {
	if !m.IsLiving() {
		return
	}

	m.Velocity = m.Velocity.Add(m.Acceleration)
	m.Hit.Center = m.Hit.Center.Add(m.Velocity)
}

func (m *Missile) Launch(start, velocity geom.PointF) {
	m.State = StateOnStage
	m.Acceleration = geom.PointF{}
	m.Hit.Center = start
	m.Velocity = velocity
}

func (m *Missile) IsLiving() bool {
	return m.State == StateOnStage
}

func (m *Missile) HitProcess(targets []Target) {
	exploded := false
	var closestTarget Target
	var closestDistance float64 = math.Inf(1)

	for _, target := range targets {
		if !target.IsLiving() {
			continue
		}
		if !target.IsEnemy() {
			continue
		}
		if !m.Hit.IntersectsWith(target.HitCircle()) {
			dist := m.Hit.Center.Distance(target.HitCircle().Center)
			if dist < closestDistance {
				closestDistance = dist
				closestTarget = target
			}
			continue
		}

		target.Damage(m.Power)
		exploded = true
	}

	if exploded {
		m.State = StateReady
		return
	}

	if closestTarget == nil {
		// There may be no enemies at all
		return
	}

	vectorForTarget := closestTarget.HitCircle().Center.Subtract(m.Hit.Center)
	targetAngle := math.Atan2(vectorForTarget.X, vectorForTarget.Y)
	m.Acceleration = geom.PointFFromPolar(m.AccelPower, targetAngle)
}

func (m *Missile) HitCircle() geom.Circle {
	return m.Hit
}

func (m *Missile) IsEnemy() bool {
	return false
}

func (m *Missile) Damage(value int) {
	m.State = StateReady
}

func (m *Missile) Position() geom.PointF {
	return m.Hit.Center
}

func (m *Missile) Angle() float64 {
	return m.Velocity.Angle()
}

func (m *Missile) VisibleF() float64 {
	if m.State == StateOnStage {
		return 1
	}
	return 0
}

func (m *Missile) Name() string {
	return name.EquipSpaceMissile
}

type EquipUpdaterHarakiriSystem struct {
	MyShipHit   geom.Circle
	Position    geom.PointF
	Interval    int
	CurrentWait int
	Harakiris   []*HarakiriSystem
	Target      geom.PointF
	HasTarget   bool
}

func (u *EquipUpdaterHarakiriSystem) Update(equip *Equip) {
	u.Position = equip.Position

	for _, h := range u.Harakiris {
		h.Update()
	}

	if u.CurrentWait > 0 {
		u.CurrentWait--
		return
	}

	u.CurrentWait = u.Interval
	for _, h := range u.Harakiris {
		if h.State == StateOnStage {
			continue
		}

		u.launchHarakiri(equip, h)
		break
	}
}

func (u *EquipUpdaterHarakiriSystem) launchHarakiri(equip *Equip, harakiri *HarakiriSystem) {
	accAngle := u.MyShipHit.Center.Subtract(equip.Position).Angle()
	var velAngle float64
	if u.HasTarget {
		velAngle = u.Target.Subtract(equip.Position).Angle()
	} else {
		velAngle = accAngle * -1
	}

	harakiri.Launch(equip.Position, velAngle, accAngle)
}

func (u *EquipUpdaterHarakiriSystem) IsLiving() bool {
	return true
}

func (u *EquipUpdaterHarakiriSystem) HitProcess(targets []Target) {
	// just aiming
	var found bool
	var closestAngle float64 = math.Pi * 2
	var closedTarget geom.PointF

	for _, target := range targets {
		if !target.IsLiving() {
			continue
		}
		if !target.IsEnemy() {
			continue
		}

		targetAngle := target.HitCircle().Center.Subtract(u.MyShipHit.Center).Angle()
		equipAngle := u.Position.Subtract(u.MyShipHit.Center).Angle()
		angleGap := math.Abs(targetAngle - equipAngle)
		if angleGap < closestAngle {
			closestAngle = angleGap
			closedTarget = target.HitCircle().Center
			found = true
		}
	}

	u.Target = closedTarget
	u.HasTarget = found
}

func (u *EquipUpdaterHarakiriSystem) Bullets() []Bullet {
	bb := make([]Bullet, len(u.Harakiris)+1)
	for i := range u.Harakiris {
		bb[i] = u.Harakiris[i]
	}
	bb[len(u.Harakiris)] = u
	return bb
}

func (u *EquipUpdaterHarakiriSystem) Targets() []Target {
	return nil
}

func (u *EquipUpdaterHarakiriSystem) VisibleEntities() []VisibleEntity {
	vv := make([]VisibleEntity, len(u.Harakiris)+1)
	for i := range u.Harakiris {
		vv[i] = u.Harakiris[i]
	}
	return vv
}

type HarakiriSystem struct {
	Hit          geom.Circle
	Velocity     geom.PointF
	FirstSpeed   float64
	Acceleration geom.PointF
	AccelPower   float64
	State        State
	Power        int
}

func (h *HarakiriSystem) Update() {
	if !h.IsLiving() {
		return
	}

	h.Velocity = h.Velocity.Add(h.Acceleration)
	h.Hit.Center = h.Hit.Center.Add(h.Velocity)
}

func (h *HarakiriSystem) Launch(start geom.PointF, valocityAngle float64, accelerationAngle float64) {
	h.State = StateOnStage
	h.Hit.Center = start
	h.Velocity = geom.PointFFromPolar(h.FirstSpeed, valocityAngle)
	h.Acceleration = geom.PointFFromPolar(h.AccelPower, accelerationAngle)
}

func (h *HarakiriSystem) IsLiving() bool {
	return h.State == StateOnStage
}

func (h *HarakiriSystem) HitProcess(targets []Target) {
	canHitMyShip := h.canHitMyShip()

	exploded := false
	for _, target := range targets {
		if !target.IsLiving() {
			continue
		}
		if !canHitMyShip && !target.IsEnemy() {
			continue
		}
		if !h.Hit.IntersectsWith(target.HitCircle()) {
			continue
		}

		target.Damage(h.Power)
		if !target.IsEnemy() {
			exploded = true
		}
	}

	if exploded {
		h.State = StateReady
	}
}

func (h *HarakiriSystem) canHitMyShip() bool {
	abs1 := h.Velocity.Abs()
	abs2 := h.Velocity.Add(h.Acceleration).Abs()

	// abs is increased, returning to myship
	return abs2 > abs1
}

func (h *HarakiriSystem) Position() geom.PointF {
	return h.Hit.Center
}

func (h *HarakiriSystem) Angle() float64 {
	return h.Velocity.Angle()
}

func (h *HarakiriSystem) VisibleF() float64 {
	if h.State == StateOnStage {
		return 1
	}
	return 0
}

func (h *HarakiriSystem) Name() string {
	return name.EquipHarakiriSystem
}

type EquipUpdaterBarrier struct {
	Hit               geom.Circle
	Durability        int
	CurrentDurability int
	Interval          int
	CurrentWait       int
}

func (u *EquipUpdaterBarrier) Update(equip *Equip) {
	u.Hit.Center = equip.Position

	if u.CurrentWait > 0 {
		u.CurrentWait--
		return
	}

	if u.CurrentDurability <= 0 {
		u.CurrentDurability = u.Durability
	}
}

func (u *EquipUpdaterBarrier) Bullets() []Bullet {
	return nil
}

func (u *EquipUpdaterBarrier) Targets() []Target {
	return []Target{u}
}

func (u *EquipUpdaterBarrier) VisibleEntities() []VisibleEntity {
	return []VisibleEntity{u}
}

func (u *EquipUpdaterBarrier) HitCircle() geom.Circle {
	return u.Hit
}

func (u *EquipUpdaterBarrier) IsEnemy() bool {
	return false
}

func (u *EquipUpdaterBarrier) Damage(value int) {
	u.CurrentDurability--
	if u.CurrentDurability <= 0 {
		u.CurrentWait = u.Interval
	}
}

func (u *EquipUpdaterBarrier) IsLiving() bool {
	return u.CurrentDurability > 0
}

func (u *EquipUpdaterBarrier) Position() geom.PointF {
	return u.Hit.Center
}

func (u *EquipUpdaterBarrier) Angle() float64 {
	return 0
}

func (u *EquipUpdaterBarrier) VisibleF() float64 {
	return float64(u.CurrentDurability) / float64(u.Durability)
}

func (u *EquipUpdaterBarrier) Name() string {
	return name.EquipBarrier
}

type EquipUpdaterExhaust struct {
	Myship     *MyShip
	Hit        geom.Circle
	Multiplier float64
}

func (u *EquipUpdaterExhaust) Update(equip *Equip) {
	u.Hit.Center = equip.Position
}

func (u *EquipUpdaterExhaust) Bullets() []Bullet {
	return nil
}

func (u *EquipUpdaterExhaust) Targets() []Target {
	return []Target{u}
}

func (u *EquipUpdaterExhaust) VisibleEntities() []VisibleEntity {
	return nil
}

func (u *EquipUpdaterExhaust) HitCircle() geom.Circle {
	return u.Hit
}

func (u *EquipUpdaterExhaust) IsEnemy() bool {
	return false
}

func (u *EquipUpdaterExhaust) Damage(value int) {
	u.Myship.Damage(int(float64(value) * u.Multiplier))
}

func (u *EquipUpdaterExhaust) IsLiving() bool {
	return u.Myship.IsLiving()
}
