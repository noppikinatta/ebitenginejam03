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

	centre := d.ShipHit.Center
	centre = centre.Add(d.StagePos)
	posTest := entity.Position()
	posTest = posTest.Add(d.StagePos)
	line := geom.LinearFuncFromPt(centre, posTest)
	x, ok := line.X(0)
	if ok {
		vector.StrokeLine(screen, float32(posTest.X), float32(posTest.Y), float32(x), 0, 2, color.RGBA{R: 255, A: 128}, false)
	}

	vector.StrokeCircle(screen, float32(posTest.X), float32(posTest.Y), 4, 1, color.RGBA{R: 255, A: 128}, false)
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

	// if m, ok := entity.(*shooter.Missile); ok {
	// fmt.Println("living missile:", *m)
	// }

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

type harakiriDrawer struct{}

func (d *harakiriDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}

type barrierDrawer struct{}

func (d *barrierDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}
