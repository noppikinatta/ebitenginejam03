package battle

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/random"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

type battleGameScene struct {
	initialized     bool
	orderer         func() []*nego.Equip
	Stage           *shooter.Stage
	VisibleEntities []shooter.VisibleEntity
	EntityDrawers   map[string]visibleEntityDrawer
	explosionDrawer *explosionDrawer
	StagePos        geom.PointF
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
	s.explosionDrawer.Update()

	return nil
}

func (s *battleGameScene) init() {
	orders := s.orderer()
	s.Stage.MyShip = s.buildMyShip(orders)
	s.Stage.EnemyLauncher = s.createEnemies()
	s.Stage.HitTest = s.createHitTest()
	s.VisibleEntities = s.createVisibleEntities()
	s.EntityDrawers = map[string]visibleEntityDrawer{
		name.TextKeyEquip1Laser: &laserDrawer{
			ShipHit:     s.Stage.MyShip.Hit,
			StagePos:    s.StagePos,
			LaserLength: (s.Stage.Size.X*0.5 - s.Stage.MyShip.Hit.Radius) * 2,
			Width:       24,
		},
		name.TextKeyEquip2Missile: &missileDrawer{
			StagePos: s.StagePos,
			Size:     8,
		},
		name.TextKeyEquip3Harakiri: &harakiriDrawer{
			StagePos: s.StagePos,
			Radius:   16,
		},
		name.TextKeyEquip4Barrier: &barrierDrawer{
			StagePos: s.StagePos,
			Radius:   48,
		},
	}
}

func (s *battleGameScene) buildMyShip(orders []*nego.Equip) *shooter.MyShip {
	myship := shooter.MyShip{
		HP:  10000,
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
			State:            shooter.EnemyStateReady,
			Hit:              geom.Circle{Radius: 8},
			ShootingInterval: 180,
			MinCloseToShip:   50,
			Bullets:          s.createBullets(),
			Rnd:              rand.New(random.Source()),
		}
	}

	return &shooter.EnemyLauncher{
		Enemies:   ee,
		Speed:     0.5,
		FirstWait: 60,
		Rnd:       rand.New(random.Source()),
		StageSize: s.Stage.Size,
		ShipHit:   s.Stage.MyShip.Hit,
		Interval:  60,
	}
}

func (s *battleGameScene) createBullets() []*shooter.EnemyBullet {
	bb := make([]*shooter.EnemyBullet, 4)

	for i := range len(bb) {
		bb[i] = &shooter.EnemyBullet{
			Power: 10,
			Hit:   geom.Circle{Radius: 1},
		}
	}

	return bb
}

func (s *battleGameScene) createHitTest() *shooter.HitTest {
	bb := s.Stage.EnemyLauncher.Bullets()
	tt := s.Stage.EnemyLauncher.Targets()
	bb = append(bb, s.Stage.MyShip.Bullets()...)
	tt = append(tt, s.Stage.MyShip.EquipTargets()...)

	exDrawer := explosionDrawer{
		StagePos: s.StagePos,
		Color:    color.RGBA{R: 255, G: 120},
	}

	for i := range bb {
		b := &explosionBullet{bullet: bb[i], Drawer: &exDrawer}
		bb[i] = b
	}

	for i := range tt {
		t := &explosionTarget{target: tt[i], Drawer: &exDrawer}
		tt[i] = t
	}

	tt = append(tt, s.Stage.MyShip)

	s.explosionDrawer = &exDrawer

	return &shooter.HitTest{
		Bullets: bb,
		Targets: tt,
	}
}

func (s *battleGameScene) createVisibleEntities() []shooter.VisibleEntity {
	return s.Stage.MyShip.VisibleEntities()
}

func (s *battleGameScene) Draw(screen *ebiten.Image) {
	s.drawBackground(screen)
	s.drawShipHP(screen)
	s.drawMyShip(screen)
	s.drawVisibleEntities(screen)
	s.drawEnemies(screen)
	s.drawExplosions(screen)
}

func (s *battleGameScene) drawRect(screen *ebiten.Image, topLeft, bottomRight geom.PointF, colorVert ebiten.Vertex) {
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	v := colorVert
	vertices := make([]ebiten.Vertex, 4)
	v.DstX = float32(topLeft.X)
	v.DstY = float32(topLeft.Y)
	vertices[0] = v
	v.DstX = float32(topLeft.X)
	v.DstY = float32(bottomRight.Y)
	vertices[1] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(bottomRight.Y)
	vertices[2] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(topLeft.Y)
	vertices[3] = v

	opt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &opt)
}

func (s *battleGameScene) drawCircle(screen *ebiten.Image, circle geom.Circle, fillColor, edgeColor color.Color) {
	cx := float32(circle.Center.X)
	cy := float32(circle.Center.Y)
	cr := float32(circle.Radius)

	vector.DrawFilledCircle(screen, cx, cy, cr, fillColor, true)
	vector.StrokeCircle(screen, cx, cy, cr, 2, edgeColor, true)
}

func (s *battleGameScene) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 48})

	v := ebiten.Vertex{
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0.5,
	}

	topLeft := geom.PointF{
		X: s.StagePos.X,
		Y: s.StagePos.Y,
	}

	bottomRight := geom.PointF{
		X: s.StagePos.X + s.Stage.Size.X,
		Y: s.StagePos.Y + s.Stage.Size.Y,
	}

	s.drawRect(screen, topLeft, bottomRight, v)
}

func (s *battleGameScene) drawShipHP(screen *ebiten.Image) {
	drawing.DrawText(screen, fmt.Sprint(s.Stage.MyShip.HP), 18, &ebiten.DrawImageOptions{})
}

func (s *battleGameScene) drawMyShip(screen *ebiten.Image) {
	circle := s.Stage.MyShip.Hit
	circle.Center = circle.Center.Add(s.StagePos)

	s.drawCircle(screen, circle, color.Gray{Y: 128}, color.Gray{Y: 96})
	for _, e := range s.Stage.MyShip.Equips {
		s.drawEquip(screen, e)
	}
}

func (s *battleGameScene) drawEquip(screen *ebiten.Image, equip *shooter.Equip) {
	center := equip.Position

	v := ebiten.Vertex{
		ColorR: 0,
		ColorG: 0.4,
		ColorB: 0.5,
		ColorA: 0.5,
	}

	topLeft := geom.PointF{
		X: center.X - 32,
		Y: center.Y - 32,
	}
	topLeft = topLeft.Add(s.StagePos)

	bottomRight := geom.PointF{
		X: center.X + 32,
		Y: center.Y + 32,
	}
	bottomRight = bottomRight.Add(s.StagePos)

	s.drawRect(screen, topLeft, bottomRight, v)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(topLeft.X, topLeft.Y)

	drawing.DrawText(screen, fmt.Sprint(equip.Name, int(equip.Angle/math.Pi*180)), 12, &opt)
}

func (s *battleGameScene) drawVisibleEntities(screen *ebiten.Image) {
	for _, e := range s.VisibleEntities {
		s.drawVisibleEntity(screen, e)
	}
}

func (s *battleGameScene) drawVisibleEntity(screen *ebiten.Image, e shooter.VisibleEntity) {
	drawer, ok := s.EntityDrawers[e.Name()]
	if !ok {
		ebitenutil.DebugPrint(screen, fmt.Sprint("WTF drawe is missing:", e.Name()))
	}
	drawer.Draw(screen, e)
}

func (s *battleGameScene) drawEnemies(screen *ebiten.Image) {
	for _, e := range s.Stage.EnemyLauncher.Enemies {
		s.drawEnemy(screen, e)
	}
}

func (s *battleGameScene) drawEnemy(screen *ebiten.Image, e *shooter.Enemy) {
	if e.State == shooter.EnemyStateOnStage {
		circle := e.Hit
		circle.Center = circle.Center.Add(s.StagePos)

		s.drawCircle(screen, circle, color.RGBA{R: 200, A: 255}, color.RGBA{R: 225, A: 255})

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(circle.Center.X, circle.Center.Y)
		drawing.DrawText(screen, fmt.Sprint(e.HP), 12, &opt)
	}

	for _, b := range e.Bullets {
		s.drawEnemyBullet(screen, b)
	}
}

func (s *battleGameScene) drawEnemyBullet(screen *ebiten.Image, b *shooter.EnemyBullet) {
	if !b.Cruising {
		return
	}

	circle := b.Hit
	circle.Center = circle.Center.Add(s.StagePos)
	s.drawCircle(screen, circle, color.RGBA{R: 255, G: 180, A: 255}, color.Transparent)
}

func (s *battleGameScene) drawExplosions(screen *ebiten.Image) {
	s.explosionDrawer.Draw(screen)
}

func (s *battleGameScene) End() bool {
	return s.Stage.End()
}

func (s *battleGameScene) Reset() {
	s.initialized = false
}
