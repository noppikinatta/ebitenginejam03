package nego

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/drawing"
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
	screen.Fill(color.Gray{Y: 48})

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
		opt := ebiten.DrawImageOptions{}

		textY := baseYOffset + 8 + (float64(i) * itemHeight)

		opt.GeoM.Translate(textX, textY)
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
