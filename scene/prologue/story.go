package prologue

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/name"
)

type Story1Drawer struct {
}

func (d *Story1Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory1)
}

type Story2Drawer struct{}

func (d *Story2Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory2)
}

type Story3Drawer struct{}

func (d *Story3Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory3)
}

type Story4Drawer struct{}

func (d *Story4Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory4)
}

func drawStoryBG(screen *ebiten.Image, key string) {
	screen.Fill(color.Gray{Y: 48})
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(24, 24)
	opt.ColorScale.Scale(1, 0.8, 0.5, 1)
	drawing.DrawTextByKey(screen, key, 24, &opt)
}
