package title

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/drawing"
)

type TitleDrawer struct {
}

func (d *TitleDrawer) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 128})
	txts := []string{
		"Totally",
		"Perfect",
		"Order",
		"to",
		"Build",
		"the",
		"Massive",
		"Space",
		"Fortress",
	}
	height := (screen.Bounds().Size().Y / len(txts)) - 1

	for i, t := range txts {
		drawStoryText(screen, float64(i*height), t)
	}
}

func drawStoryText(screen *ebiten.Image, yOffset float64, txt string) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(24, yOffset)
	opt.ColorScale.Scale(0.8, 0.8, 0.8, 1)

	drawing.DrawText(screen, txt, 48, &opt)
}
