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
	initialized bool
	orderer     func() []*nego.Equip
	Stage       *shooter.Stage
	StagePos    geom.PointF
}

func newBattleGameScene(orderer func() []*nego.Equip) *battleGameScene {
	s := shooter.Stage{
		Size: geom.PointF{X: 600, Y: 600},
	}

	return &battleGameScene{
		orderer:  orderer,
		Stage:    &s,
		StagePos: geom.PointF{X: 0, Y: 40},
	}
}

func (s *battleGameScene) Update() error {
	if !s.initialized {
		s.init()
		s.initialized = true
	}

	x, y := ebiten.CursorPosition()
	cursorPos := geom.PointF{X: float64(x), Y: float64(y)}
	cursorPos = cursorPos.Subtract(s.StagePos)
	s.Stage.Update(cursorPos)

	return nil
}

func (s *battleGameScene) init() {
	orders := s.orderer()
	s.Stage.MyShip = s.buildMyShip(orders)
	s.Stage.EnemyLauncher = s.createEnemies()
	s.Stage.HitTest = s.createHitTest()
}

func (s *battleGameScene) buildMyShip(orders []*nego.Equip) *shooter.MyShip {
	myship := shooter.MyShip{
		HP:  1000,
		Hit: geom.Circle{Center: s.Stage.Size.Multiply(0.5), Radius: 100},
	}

	build.BuildEquips(&myship, orders)

	return &myship
}

func (s *battleGameScene) createEnemies() *shooter.EnemyLauncher {
	ee := make([]*shooter.Enemy, 100)

	for i := range len(ee) {
		ee[i] = &shooter.Enemy{
			HP:               100,
			State:            shooter.StateReady,
			Hit:              geom.Circle{Radius: 32},
			ShootingInterval: 180,
			Bullets:          s.createBullets(),
			Rnd:              rand.New(random.Source()),
		}
	}

	return &shooter.EnemyLauncher{
		Enemies:   ee,
		Speed:     1,
		FirstWait: 180,
		Rnd:       rand.New(random.Source()),
		StageSize: s.Stage.Size,
		Interval:  60,
	}
}

func (s *battleGameScene) createBullets() []*shooter.EnemyBullet {
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

func (s *battleGameScene) createHitTest() *shooter.HitTest {
	bb := s.Stage.EnemyLauncher.Bullets()
	tt := s.Stage.EnemyLauncher.Targets()
	bb = append(bb, s.Stage.MyShip.Bullets()...)
	tt = append(tt, s.Stage.MyShip.Targets()...)

	return &shooter.HitTest{
		Bullets: bb,
		Targets: tt,
	}
}

func (s *battleGameScene) Draw(screen *ebiten.Image) {

}

func (s *battleGameScene) End() bool {
	return s.Stage.End()
}

func (s *battleGameScene) Reset() {
	s.initialized = false
}
