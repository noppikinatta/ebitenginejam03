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
	"github.com/noppikinatta/ebitenginejam03/scene"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

type battleGameScene struct {
	initialized     bool
	orderer         func() []*nego.Equip
	Stage           *shooter.Stage
	VisibleEntities []shooter.VisibleEntity
	EntityDrawers   map[string]visibleEntityDrawer
	explosionDrawer *explosionDrawer
	hpDrawer        *drawing.GaugeDrawer
	StagePos        geom.PointF
	preprocess      scene.Scene
	postprocess     scene.Scene
	bgImg           *ebiten.Image
}

func newBattleGameScene(orderer func() []*nego.Equip) *battleGameScene {
	s := shooter.Stage{
		Size: geom.PointF{X: 600, Y: 600},
	}

	cmax := ebiten.ColorScale{}
	cmax.Scale(0.75, 0.75, 0.75, 1)

	cmin := ebiten.ColorScale{}
	cmin.SetG(0.25)
	cmin.SetB(0.25)

	endScene := battleEndScene{
		resultFn: s.Won,
		frames:   60,
	}

	return &battleGameScene{
		orderer:  orderer,
		Stage:    &s,
		StagePos: geom.PointF{X: 0, Y: 40},
		hpDrawer: &drawing.GaugeDrawer{
			TopLeft:       geom.PointF{X: 0, Y: 0},
			BottomRight:   geom.PointF{X: 600, Y: 40},
			TextOffset:    geom.PointF{X: 64, Y: 6},
			FontSize:      18,
			ColorScaleMax: cmax,
			ColorScaleMin: cmin,
		},
		preprocess: scene.NewContainer(
			scene.NewFadeIn(15),
			scene.NewShowImageScene(0, &ruleDescriptionDrawer{}),
			&scene.ScrollText{TextKey: name.TextKeyShooterTitle2},
		),
		postprocess: scene.NewContainer(
			&endScene,
			scene.NewShowImageScene(0, &endScene),
			scene.NewFadeOut(30),
		),
		bgImg: CreateBG(int(s.Size.X), int(s.Size.Y)),
	}
}

func (s *battleGameScene) Update() error {
	if !s.initialized {
		s.init()
		s.initialized = true
	}

	if !s.Stage.End() {
		x, y := ebiten.CursorPosition()
		cursorPos := geom.PointF{X: float64(x), Y: float64(y)}
		cursorPos = cursorPos.Subtract(s.StagePos)
		s.Stage.UpdateAngle(cursorPos)
	}
	if !s.preprocess.End() {
		return s.preprocess.Update()
	}
	if s.Stage.End() {
		s.updateMyshipExplosion()
		return s.postprocess.Update()
	}
	s.Stage.UpdateOther()
	s.explosionDrawer.Update()
	s.hpDrawer.Current = s.Stage.MyShip.HP

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
	s.hpDrawer.Max = s.Stage.MyShip.HP
	s.hpDrawer.Current = s.Stage.MyShip.HP
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
	ee := make([]*shooter.Enemy, 80)

	for i := range len(ee) {
		ee[i] = &shooter.Enemy{
			HP:               100,
			MaxHP:            100,
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

func (s *battleGameScene) updateMyshipExplosion() {

}

func (s *battleGameScene) Draw(screen *ebiten.Image) {
	if !s.initialized {
		screen.Fill(color.Black)
		return
	}
	s.drawBackground(screen)
	s.drawMyShip(screen)
	s.drawVisibleEntities(screen)
	s.drawEnemies(screen)
	s.drawExplosions(screen)
	s.drawMyshipExplosion(screen)
	s.drawShipHP(screen)
	s.drawEnemyList(screen)
	if !s.preprocess.End() {
		s.preprocess.Draw(screen)
	}
	if s.Stage.End() && !s.postprocess.End() {
		s.postprocess.Draw(screen)
	}
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
	screen.Fill(color.Gray{Y: 24})

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(s.StagePos.X, s.StagePos.Y)
	screen.DrawImage(s.bgImg, &opt)
}

func (s *battleGameScene) drawShipHP(screen *ebiten.Image) {
	v := ebiten.Vertex{}
	gray := float32(24) / float32(255)
	v.ColorR = gray
	v.ColorG = gray
	v.ColorB = gray
	v.ColorA = 1
	s.drawRect(screen, s.hpDrawer.TopLeft, s.hpDrawer.BottomRight, v)
	s.hpDrawer.Draw(screen)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(8, s.hpDrawer.TextOffset.Y)
	drawing.DrawText(screen, "HP:", s.hpDrawer.FontSize, &opt)
}

func (s *battleGameScene) drawMyShip(screen *ebiten.Image) {
	circle := s.Stage.MyShip.Hit
	circle.Center = circle.Center.Add(s.StagePos)

	s.drawCircle(screen, circle, color.Gray{Y: 96}, color.Gray{Y: 48})

	shipImg := drawing.Image(name.ImgKeyMyship)
	shipImgSize := geom.PointFFromPoint(shipImg.Bounds().Size())
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(-shipImgSize.X*0.5, -shipImgSize.Y*0.5)
	opt.GeoM.Translate(circle.Center.X, circle.Center.Y)
	screen.DrawImage(shipImg, &opt)

	for _, e := range s.Stage.MyShip.Equips {
		s.drawEquip(screen, e)
	}
}

func (s *battleGameScene) drawEquip(screen *ebiten.Image, equip *shooter.Equip) {
	if equip.Name == name.TextKeyEquip4Barrier {
		return
	}

	center := equip.Position.Add(s.StagePos)

	eqpImg := drawing.Image(name.ImgKey(equip.Name))
	eqpImgSize := geom.PointFFromPoint(eqpImg.Bounds().Size())

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(-eqpImgSize.X*0.5, -eqpImgSize.Y*0.5)
	opt.GeoM.Scale(0.5, 0.5)
	opt.GeoM.Rotate(equip.Angle + 0.5*math.Pi)
	if equip.Name == name.TextKeyEquip1Laser {
		opt.GeoM.Rotate(0.5 * math.Pi)
	}
	opt.GeoM.Translate(center.X, center.Y)
	if upd, ok := equip.Updater.(*shooter.EquipUpdaterExhaust); ok {
		if upd.Alpha > 0 {
			gb := float32(1 - upd.Alpha)
			opt.ColorScale.Scale(1, gb, gb, 1)
		}
	}
	screen.DrawImage(eqpImg, &opt)
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
		enemyImg := drawing.Image(name.ImgKeyEnemy)
		enemyImgSize := geom.PointFFromPoint(enemyImg.Bounds().Size())

		circle := e.Hit
		circle.Center = circle.Center.Add(s.StagePos)

		opt := ebiten.DrawImageOptions{}
		r := s.Stage.MyShip.Hit.Center.Subtract(e.Hit.Center).Angle()
		opt.GeoM.Translate(-enemyImgSize.X*0.5, -enemyImgSize.Y*0.5)
		opt.GeoM.Rotate(r)
		opt.GeoM.Translate(circle.Center.X, circle.Center.Y)
		hpp := float32(e.HP) / float32(e.MaxHP)
		opt.ColorScale.Scale(1, hpp, hpp, 1)

		screen.DrawImage(enemyImg, &opt)
	}

	for _, b := range e.Bullets {
		s.drawEnemyBullet(screen, b)
	}
}

func (s *battleGameScene) drawEnemyList(screen *ebiten.Image) {
	const (
		leftOffset float64 = 600
		clmSize    float64 = 40
		rowSize    float64 = 40
		clmCount           = 5
	)

	v := ebiten.Vertex{}
	gray := float32(24) / float32(255)
	v.ColorR = gray
	v.ColorG = gray
	v.ColorB = gray
	v.ColorA = 1
	screenSize := geom.PointFFromPoint(screen.Bounds().Size())
	s.drawRect(screen, geom.PointF{X: leftOffset}, screenSize, v)

	enemyImg := drawing.Image(name.ImgKeyEnemy)
	enemyImgSize := geom.PointFFromPoint(enemyImg.Bounds().Size())

	for i, e := range s.Stage.EnemyLauncher.Enemies {
		if e.State == shooter.EnemyStateDead {
			continue
		}

		x := float64(i%clmCount) * clmSize
		x += leftOffset + (clmSize-enemyImgSize.X)*0.5
		y := float64(i/clmCount) * rowSize
		y += (rowSize - enemyImgSize.Y) * 0.5

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)

		if e.State == shooter.EnemyStateOnStage {
			hpp := float32(e.HP) / float32(e.MaxHP)
			opt.ColorScale.Scale(1, hpp, hpp, 1)
		} else {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 0.5)
		}

		screen.DrawImage(enemyImg, &opt)
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

func (s *battleGameScene) drawMyshipExplosion(screen *ebiten.Image) {

}

func (s *battleGameScene) End() bool {
	return s.postprocess.End()
}

func (s *battleGameScene) Reset() {
	s.preprocess.Reset()
	s.postprocess.Reset()
	s.initialized = false
}

type ruleDescriptionDrawer struct {
}

func (d *ruleDescriptionDrawer) Draw(screen *ebiten.Image) {
	size := screen.Bounds().Size()
	clr := color.RGBA{A: 128}
	vector.DrawFilledRect(screen, 0, 0, float32(size.X), float32(size.Y), clr, false)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(24, 280)
	drawing.DrawTextByKey(screen, name.TextKeyShooterDesc1, 18, &opt)
}
