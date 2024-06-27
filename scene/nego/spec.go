package nego

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
)

type specDrawer struct {
	Orderer   func() []*nego.Equip
	equipDesc *build.EquipDescriptor
}

func newSpecDrawer(orderer func() []*nego.Equip) *specDrawer {
	return &specDrawer{
		Orderer:   orderer,
		equipDesc: build.NewEquipDescriptor(),
	}
}

func (d *specDrawer) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 96})

	d.drawTitle(screen)
	d.drawEquips(screen)
}

func (d *specDrawer) drawTitle(screen *ebiten.Image) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(4, 4)

	drawing.DrawTextByKey(screen, name.TextKeyNegotiationTitle3, 20, &opt)
}

func (d *specDrawer) drawEquips(screen *ebiten.Image) {
	equips := d.Orderer()
	if len(equips) == 0 {
		return
	}

	itemHeight := math.Min(80, 600/float64(len(equips)))

	const (
		baseXOffset  float64 = 4
		baseYOffset  float64 = 40
		imageXOffset float64 = 64
	)

	textX := baseXOffset*2 + imageXOffset

	for i, e := range equips {
		if i%2 == 0 {
			tl := geom.PointF{
				X: 0,
				Y: baseYOffset + (float64(i) * itemHeight),
			}
			br := geom.PointF{
				X: float64(screen.Bounds().Max.X),
				Y: baseYOffset + (float64(i+1) * itemHeight),
			}

			cv := ebiten.Vertex{
				ColorA: 0.5,
			}

			d.drawRect(screen, tl, br, cv)
		}

		img := drawing.Image(name.ImgKey(e.Name))
		imgSize := geom.PointFFromPoint(img.Bounds().Size())

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(baseXOffset, baseYOffset+(float64(i)*itemHeight))

		iopt := opt
		iopt.GeoM.Translate(0, (itemHeight-imgSize.Y)*0.5)
		screen.DrawImage(img, &iopt)

		opt.GeoM.Translate(textX, 8)
		eqpName := lang.Text(e.Name)
		if e.ImprovedCount > 0 {
			eqpName += fmt.Sprintf("+%d", e.ImprovedCount)
		}
		drawing.DrawText(screen, eqpName, 14, &opt)

		opt.GeoM.Translate(0, 20)

		eqpDesc := name.DescKey(e.Name)
		tmplData := d.equipDesc.TemplateData(e.Name, e.ImprovedCount)
		drawing.DrawTextTemplate(screen, eqpDesc, tmplData, 12, &opt)
	}
}

func (d *specDrawer) drawRect(screen *ebiten.Image, topLeft, bottomRight geom.PointF, colorVert ebiten.Vertex) {
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
