package battle

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

type visibleEntityDrawer interface {
	Draw(screen *ebiten.Image, entity shooter.VisibleEntity)
}

type laserDrawer struct {
	ShipHit     geom.Circle
	StagePos    geom.PointF
	LaserLength float64
	Width       float64
}

func (d *laserDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {
	colorVert := ebiten.Vertex{
		ColorR: 0.9,
		ColorG: 0.95,
		ColorB: 1,
		ColorA: d.alpha(entity),
	}

	topLeft := geom.PointF{
		X: 0,
		Y: d.Width * -0.5,
	}

	bottomRight := geom.PointF{
		X: d.LaserLength,
		Y: d.Width * 0.5,
	}

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

	gm := ebiten.GeoM{}
	gm.Translate(d.ShipHit.Radius, 0)
	gm.Rotate(entity.Angle())
	gm.Translate(d.ShipHit.Center.X, d.ShipHit.Center.Y)
	gm.Translate(d.StagePos.X, d.StagePos.Y)

	for i := range vertices {
		vt := vertices[i]
		x, y := vt.DstX, vt.DstY
		x64, y64 := gm.Apply(float64(x), float64(y))
		vt.DstX = float32(x64)
		vt.DstY = float32(y64)
		vertices[i] = vt
	}

	opt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &opt)
}

func (d *laserDrawer) alpha(entity shooter.VisibleEntity) float32 {
	if entity.VisibleF() > 0.2 {
		return 1
	}

	return float32(entity.VisibleF() * 5)
}

type missileDrawer struct {
	StagePos geom.PointF
	Size     float64
}

func (d *missileDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {
	if entity.VisibleF() == 0 {
		return
	}

	pos := entity.Position()

	gm := ebiten.GeoM{}
	gm.Rotate(entity.Angle())
	gm.Translate(pos.X, pos.Y)
	gm.Translate(d.StagePos.X, d.StagePos.Y)
	var x, y float64
	v := ebiten.Vertex{
		ColorR: 0.2,
		ColorG: 0.8,
		ColorB: 0.6,
		ColorA: 1,
	}

	vertices := make([]ebiten.Vertex, 3)
	x, y = gm.Apply(d.Size*0.5, 0)
	v.DstX = float32(x)
	v.DstY = float32(y)
	vertices[0] = v
	x, y = gm.Apply(-d.Size*0.5, -d.Size*0.25)
	v.DstX = float32(x)
	v.DstY = float32(y)
	vertices[1] = v
	x, y = gm.Apply(-d.Size*0.5, d.Size*0.5)
	v.DstX = float32(x)
	v.DstY = float32(y)
	vertices[2] = v

	topt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, []uint16{0, 1, 2}, drawing.WhitePixel, &topt)
}

type harakiriDrawer struct {
	StagePos geom.PointF
	Radius   float32
}

func (d *harakiriDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {
	if entity.VisibleF() == 0 {
		return
	}

	pos := entity.Position().Add(d.StagePos)

	cx := float32(pos.X)
	cy := float32(pos.Y)

	vector.DrawFilledCircle(screen, cx, cy, d.Radius, color.RGBA{R: 100, G: 200, A: 255}, true)
}

type barrierDrawer struct {
	StagePos geom.PointF
	Radius   float32
}

func (d *barrierDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {
	if entity.VisibleF() == 0 {
		return
	}

	pos := entity.Position().Add(d.StagePos)

	cx := float32(pos.X)
	cy := float32(pos.Y)

	var r, g, b, a uint8 = 128, 128, 255, 255

	redModifier := 1 - entity.VisibleF()
	r = uint8(120*redModifier) + r

	b = uint8(float64(b) * entity.VisibleF())

	fillAlpha := entity.VisibleF() * 0.5

	f := func(v uint8) uint8 {
		return uint8(float64(v) * fillAlpha)
	}

	edgeColor := color.RGBA{R: r, G: g, B: b, A: a}
	fillColor := color.RGBA{R: f(r), G: f(g), B: f(b), A: f(a)}

	vector.DrawFilledCircle(screen, cx, cy, d.Radius, fillColor, true)
	vector.StrokeCircle(screen, cx, cy, d.Radius, 2, edgeColor, true)
}

type explosionDrawer struct {
	Color      color.Color
	StagePos   geom.PointF
	explosions []*explosion
}

func (d *explosionDrawer) Add(at geom.Circle) {
	for _, e := range d.explosions {
		if e.Alpha == 0 {
			e.At = at
			e.Alpha = 1
			return
		}
	}

	d.explosions = append(d.explosions, &explosion{At: at, Alpha: 1})
}

func (d *explosionDrawer) Update() {
	for _, e := range d.explosions {
		e.Update()
	}
}

func (d *explosionDrawer) Draw(screen *ebiten.Image) {
	for _, e := range d.explosions {
		e.Draw(screen, d.StagePos, d.Color)
	}
}

type explosion struct {
	At    geom.Circle
	Alpha float64
}

func (e *explosion) Update() {
	if e.Alpha <= 0 {
		return
	}

	e.Alpha -= 0.01
	if e.Alpha < 0 {
		e.Alpha = 0
	}
}

func (e *explosion) Draw(screen *ebiten.Image, stagePos geom.PointF, clr color.Color) {
	if e.Alpha == 0 {
		return
	}

	cx := float32(e.At.Center.X + stagePos.X)
	cy := float32(e.At.Center.Y + stagePos.Y)
	cr := float32(e.At.Radius)

	r, g, b, _ := clr.RGBA()
	r16 := uint16(float64(r) * e.Alpha)
	g16 := uint16(float64(g) * e.Alpha)
	b16 := uint16(float64(b) * e.Alpha)
	a16 := uint16(float64(0x00ff) * e.Alpha)
	c := color.RGBA64{R: r16, G: g16, B: b16, A: a16}

	vector.DrawFilledCircle(screen, cx, cy, cr, c, true)
}

type explosionBullet struct {
	bullet shooter.Bullet
	Drawer *explosionDrawer
}

func (b *explosionBullet) IsLiving() bool {
	return b.bullet.IsLiving()
}

func (b *explosionBullet) HitProcess(targets []shooter.Target) geom.Circle {
	explosionCircle := b.bullet.HitProcess(targets)
	if explosionCircle.Radius > 0 {
		b.Drawer.Add(explosionCircle)
	}

	return explosionCircle
}

type explosionTarget struct {
	target shooter.Target
	Drawer *explosionDrawer
}

func (t *explosionTarget) HitCircle() geom.Circle {
	return t.target.HitCircle()
}

func (t *explosionTarget) IsEnemy() bool {
	return t.target.IsEnemy()
}

func (t *explosionTarget) Damage(value int) float64 {
	explosionRadius := t.target.Damage(value)
	if explosionRadius > 0 {
		at := t.target.HitCircle()
		at.Radius = explosionRadius
		t.Drawer.Add(at)
	}
	return explosionRadius
}

func (t *explosionTarget) IsLiving() bool {
	return t.target.IsLiving()
}
