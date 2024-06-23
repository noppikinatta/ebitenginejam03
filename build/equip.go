package build

import (
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

func createBuilders() map[string]builder {
	return map[string]builder{
		name.EquipLaserCannon: &laserBuilder{
			LastFrames: 200,
			Interval:   300,
			Width:      24,
			Power:      10,
		},
		name.EquipSpaceMissile: &missileBuilder{
			Interval:   300,
			MaxCount:   6,
			HitRadius:  32,
			FirstSpeed: 10,
			AccelPower: 2,
			Power:      50,
		},
		name.EquipHarakiriSystem: &harakiriSystemBuilder{
			Interval:   180,
			MaxCount:   6,
			HitRadius:  32,
			FirstSpeed: 20,
			AccelPower: 2,
			Power:      100,
		},
		name.EquipBarrier: &barrierBuilder{
			HitRadius:  64,
			Durability: 5,
			Interval:   240,
		},
		name.EquipArmorPlate: &armorPlatebuilder{
			AdditinalArmor: 100,
		},
		name.EquipThermalExhaustPort: &exhaustPortbuilder{
			HitRadius:  32,
			Multiplier: 10,
		},
		name.EquipStonehenge: &uselessBuilder{
			Name: name.EquipStonehenge,
		},
		name.EquipSushiBar: &uselessBuilder{
			Name: name.EquipSushiBar,
		},
		name.EquipOperaHouse: &uselessBuilder{
			Name: name.EquipOperaHouse,
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
		Name: name.EquipLaserCannon,
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
		"LastFrames": b.calcedLastFrames(improvedCount),
		"Interval":   b.Interval,
		"Power":      b.calcedPower(improvedCount),
	}
}

type missileBuilder struct {
	Interval   int
	MaxCount   int
	HitRadius  float64
	FirstSpeed float64
	AccelPower float64
	Power      int
}

func (b *missileBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.EquipSpaceMissile,
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
			Hit:        geom.Circle{Radius: b.HitRadius},
			FirstSpeed: b.FirstSpeed,
			AccelPower: b.AccelPower,
			State:      shooter.StateReady,
			Power:      b.calcedPower(improvedCount),
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
		"Interval": b.calcedInterval(improvedCount),
		"Power":    b.calcedPower(improvedCount),
	}
}

type harakiriSystemBuilder struct {
	Interval   int
	MaxCount   int
	HitRadius  float64
	FirstSpeed float64
	AccelPower float64
	Power      int
}

func (b *harakiriSystemBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.EquipHarakiriSystem,
		Updater: &shooter.EquipUpdaterHarakiriSystem{
			Interval:  b.calcedInterval(improvedCount),
			Harakiris: b.buildHarakiris(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *harakiriSystemBuilder) buildHarakiris(improvedCount int) []*shooter.HarakiriSystem {
	maxCount := b.calcedMaxCount(improvedCount)
	hh := make([]*shooter.HarakiriSystem, maxCount)

	for i := range maxCount {
		hh[i] = &shooter.HarakiriSystem{
			Hit:        geom.Circle{Radius: b.HitRadius},
			FirstSpeed: b.FirstSpeed,
			AccelPower: b.AccelPower,
			Power:      b.Power,
		}
	}

	return hh
}

func (b *harakiriSystemBuilder) calcedInterval(improvedCount int) int {
	return improve(b.Interval, 0.5, improvedCount)
}

func (b *harakiriSystemBuilder) calcedMaxCount(improvedCount int) int {
	return improve(b.MaxCount, 2, improvedCount)
}

func (b *harakiriSystemBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Interval": b.calcedInterval(improvedCount),
		"Power":    b.Power,
	}
}

type barrierBuilder struct {
	HitRadius  float64
	Durability int
	Interval   int
}

func (b *barrierBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name: name.EquipBarrier,
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
		"Durability": b.calcedDurability(improvedCount),
		"Interval":   b.calcedInterval(improvedCount),
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
		Name: name.EquipThermalExhaustPort,
		Updater: &shooter.EquipUpdaterExhaust{
			Myship:     ship,
			Hit:        geom.Circle{Radius: b.HitRadius},
			Multiplier: b.calcedMultiplier(improvedCount),
		},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *exhaustPortbuilder) calcedMultiplier(improvedCount int) float64 {
	return improve(b.Multiplier, 1.5, improvedCount)
}

func (b *exhaustPortbuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{
		"Multiplier": b.calcedMultiplier(improvedCount),
	}
}

type uselessBuilder struct {
	Name string
}

func (b *uselessBuilder) Build(ship *shooter.MyShip, improvedCount int) {
	eqp := shooter.Equip{
		Name:    b.Name,
		Updater: &shooter.EquipUpdaterNop{},
	}

	ship.Equips = append(ship.Equips, &eqp)
}

func (b *uselessBuilder) TemplateData(improvedCount int) map[string]any {
	return map[string]any{}
}

func improve[T int | float64](baseValue T, rate float64, count int) T {
	v := float64(baseValue)

	for range count {
		v *= rate
	}

	return T(v)
}
