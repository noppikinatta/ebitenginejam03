package build

import (
	"fmt"

	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

func BuildEquips(baseShip *shooter.MyShip, orders []*nego.Equip) {
	builders := createBuilders()

	for _, o := range orders {
		builder := builders[o.Name]
		builder.Build(baseShip, o.ImprovedCount)
	}
}

type EquipDescriptor struct {
	builders map[string]builder
}

func NewEquipDescriptor() *EquipDescriptor {
	d := EquipDescriptor{
		builders: createBuilders(),
	}
	return &d
}

func (d *EquipDescriptor) TemplateData(txtKey string, improvedCount int) map[string]any {
	b := d.builders[txtKey]
	return b.TemplateData(improvedCount)
}

func createBuilders() map[string]builder {
	return map[string]builder{
		name.TextKeyEquip1Laser: &laserBuilder{
			LastFrames: 200,
			Interval:   300,
			Width:      24,
			Power:      10,
		},
		name.TextKeyEquip2Missile: &missileBuilder{
			Interval:       90,
			MaxCount:       6,
			HitRadius:      8,
			ExplodeRadius:  32,
			FirstSpeed:     0.5,
			AccelPower:     0.125,
			Power:          50,
			LifetimeFrames: 180,
		},
		name.TextKeyEquip3Harakiri: &harakiriSystemBuilder{
			Interval:       300,
			MaxCount:       1,
			HitRadius:      16,
			FirstSpeed:     8,
			Power:          300,
			MaxSanity:      3,
			AimingInterval: 120,
		},
		name.TextKeyEquip4Barrier: &barrierBuilder{
			HitRadius:  48,
			Durability: 10,
			Interval:   240,
		},
		name.TextKeyEquip5Armor: &armorPlatebuilder{
			AdditinalArmor: 2000,
		},
		name.TextKeyEquip6Exhaust: &exhaustPortbuilder{
			HitRadius:  16,
			Multiplier: 10,
		},
		name.TextKeyEquip7Stonehenge: &uselessBuilder{
			Name:  name.TextKeyEquip7Stonehenge,
			Value: 30,
		},
		name.TextKeyEquip8Sushibar: &uselessBuilder{
			Name:  name.TextKeyEquip8Sushibar,
			Value: 150,
		},
		name.TextKeyEquip9Operahouse: &uselessBuilder{
			Name:  name.TextKeyEquip9Operahouse,
			Value: 1507,
		},
	}
}

type builder interface {
	Build(ship *shooter.MyShip, improvedCount int)
	TemplateData(improvedCount int) map[string]any
}

type laserBuilder struct {
	LastFrames int
	Interval   int
	Width      float64
	Power      int
}

func (b *laserBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.TextKeyEquip1Laser,
		Updater: &shooter.EquipUpdaterLaser{
			ShipHit:    ship.Hit,
			LastFrames: b.calcedLastFrames(improvedCount),
			Interval:   b.Interval,
			Width:      b.Width,
			Power:      b.calcedPower(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *laserBuilder) calcedLastFrames(improvedCount int) int {
	return improve(b.LastFrames, 1.25, improvedCount)
}

func (b *laserBuilder) calcedPower(improvedCount int) int {
	return improve(b.Power, 1.25, improvedCount)
}

func (b *laserBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"LastSec":     ftos(b.calcedLastFrames(improvedCount)),
		"IntervalSec": ftos(b.Interval),
		"Power":       b.calcedPower(improvedCount) * 60,
	}
}

type missileBuilder struct {
	Interval       int
	MaxCount       int
	HitRadius      float64
	ExplodeRadius  float64
	FirstSpeed     float64
	AccelPower     float64
	Power          int
	LifetimeFrames int
}

func (b *missileBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.TextKeyEquip2Missile,
		Updater: &shooter.EquipUpdaterMissile{
			Interval: b.calcedInterval(improvedCount),
			Missiles: b.buildMissiles(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *missileBuilder) buildMissiles(improvedCount int) []*shooter.Missile {
	mm := make([]*shooter.Missile, b.MaxCount)

	for i := range b.MaxCount {
		mm[i] = &shooter.Missile{
			Hit:            geom.Circle{Radius: b.HitRadius},
			FirstSpeed:     b.FirstSpeed,
			AccelPower:     b.AccelPower,
			State:          shooter.MissileStateReady,
			Power:          b.calcedPower(improvedCount),
			ExplodeRadius:  b.ExplodeRadius,
			LifetimeFrames: b.LifetimeFrames,
		}
	}

	return mm
}

func (b *missileBuilder) calcedInterval(improvedCount int) int {
	return improve(b.Interval, 0.8, improvedCount)
}

func (b *missileBuilder) calcedPower(improvedCount int) int {
	return improve(b.Power, 1.25, improvedCount)
}

func (b *missileBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"IntervalSec": ftos(b.calcedInterval(improvedCount)),
		"Power":       b.calcedPower(improvedCount),
	}
}

type harakiriSystemBuilder struct {
	Interval       int
	MaxCount       int
	HitRadius      float64
	FirstSpeed     float64
	Power          int
	MaxSanity      int
	AimingInterval int
}

func (b *harakiriSystemBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.TextKeyEquip3Harakiri,
		Updater: &shooter.EquipUpdaterHarakiriSystem{
			Interval:  b.Interval,
			MaxSanity: b.calcedMaxSanity(improvedCount),
			Harakiris: b.buildHarakiris(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *harakiriSystemBuilder) buildHarakiris(improvedCount int) []*shooter.HarakiriSystem {
	hh := make([]*shooter.HarakiriSystem, b.MaxCount)

	for i := range b.MaxCount {
		hh[i] = &shooter.HarakiriSystem{
			Hit:            geom.Circle{Radius: b.HitRadius},
			FirstSpeed:     b.FirstSpeed,
			Power:          b.calcedPower(improvedCount),
			AimingInterval: b.AimingInterval,
		}
	}

	return hh
}

func (b *harakiriSystemBuilder) calcedMaxSanity(improvedCount int) int {
	return improve(b.MaxSanity, 1.5, improvedCount)
}

func (b *harakiriSystemBuilder) calcedPower(improvedCount int) int {
	return improve(b.Power, 1.2, improvedCount)
}

func (b *harakiriSystemBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"MaxSanity": b.calcedMaxSanity(improvedCount),
		"Power":     b.calcedPower(improvedCount),
	}
}

type barrierBuilder struct {
	HitRadius  float64
	Durability int
	Interval   int
}

func (b *barrierBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.TextKeyEquip4Barrier,
		Updater: &shooter.EquipUpdaterBarrier{
			Hit:        geom.Circle{Radius: b.HitRadius},
			Durability: b.calcedDurability(improvedCount),
			Interval:   b.calcedInterval(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *barrierBuilder) calcedDurability(improvedCount int) int {
	return improve(b.Durability, 1.5, improvedCount)
}

func (b *barrierBuilder) calcedInterval(improvedCount int) int {
	return improve(b.Interval, 0.75, improvedCount)
}

func (b *barrierBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Durability":  b.calcedDurability(improvedCount),
		"IntervalSec": ftos(b.calcedInterval(improvedCount)),
	}
}

type armorPlatebuilder struct {
	AdditinalArmor int
}

func (b *armorPlatebuilder) Build(ship *shooter.MyShip, improvedCount int) {
	ship.HP += b.calcedAdditionalArmor(improvedCount)
}

func (b *armorPlatebuilder) calcedAdditionalArmor(improvedCount int) int {
	return improve(b.AdditinalArmor, 2, improvedCount)
}

func (b *armorPlatebuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Armor": b.calcedAdditionalArmor(improvedCount),
	}
}

type exhaustPortbuilder struct {
	HitRadius  float64
	Multiplier float64
}

func (b *exhaustPortbuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.TextKeyEquip6Exhaust,
		Updater: &shooter.EquipUpdaterExhaust{
			Myship:     ship,
			Hit:        geom.Circle{Radius: b.HitRadius},
			Multiplier: b.calcedMultiplier(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *exhaustPortbuilder) calcedMultiplier(improvedCount int) float64 {
	return improve(b.Multiplier, 1.2, improvedCount)
}

func (b *exhaustPortbuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Multiplier": int(b.calcedMultiplier(improvedCount)),
	}
}

type uselessBuilder struct {
	Name  string
	Value int
}

func (b *uselessBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name:    b.Name,
		Updater: &shooter.EquipUpdaterNop{},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *uselessBuilder) calcedValue(improvedCount int) int {
	return improve(b.Value, 1.2, improvedCount)
}

func (b *uselessBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Value": b.calcedValue(improvedCount),
	}
}

func improve[T int | float64](baseValue T, rate float64, count int) T {
	v := float64(baseValue)

	for range count {
		v *= rate
	}

	return T(v)
}

func ftos(frames int) string {
	s := float64(frames) / 60
	return fmt.Sprintf("%.1f", s)
}
