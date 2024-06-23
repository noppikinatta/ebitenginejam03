package battle

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/random"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

type battleGameScene struct {
	Stage    *shooter.Stage
	StagePos geom.PointF
}

func newBattleGameScene(orders []*nego.Equip) *battleGameScene {
	s := shooter.Stage{
		Size: geom.PointF{X: 600, Y: 600},
	}

	s.MyShip = buildMyShip(s.Size, orders)

	return &battleGameScene{
		Stage: &s,
	}
}

func buildMyShip(stageSize geom.PointF, orders []*nego.Equip) *shooter.MyShip {
	myship := shooter.MyShip{
		HP:  1000,
		Hit: geom.Circle{Center: stageSize.Multiply(0.5), Radius: stageSize.Abs() * 0.2},
	}

	build.BuildEquips(&myship, orders)

	return &myship
}

func createEnemies() []*shooter.Enemy {
	ee := make([]*shooter.Enemy, 100)

	for i := range len(ee) {
		ee[i] = &shooter.Enemy{
			HP:               100,
			State:            shooter.StateReady,
			Hit:              geom.Circle{Radius: 32},
			ShootingInterval: 180,
			Bullets:          createBullets(),
			Rnd:              rand.New(random.Source()),
		}
	}

	return ee
}

func createBullets() []*shooter.EnemyBullet {
	bb := make([]*shooter.EnemyBullet, 4)

	for i := range len(bb) {
		bb[i] = &shooter.EnemyBullet{
			Power: 10,
			Hit:   geom.Circle{Radius: 1},
			State: shooter.StateReady,
		}
	}

	return bb
}

func (s *battleGameScene) Update() error {

}

func (s *battleGameScene) Draw(screen *ebiten.Image) {

}

func (s *battleGameScene) End() bool {

}

func (s *battleGameScene) Reset() {

}
