package nego

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
)

type prepareDrawer struct {
}

func (d *prepareDrawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyShooterTitle1)
}

func drawStoryBG(screen *ebiten.Image, key string) {
	screen.Fill(color.Gray{Y: 24})

	img := drawing.Image(name.ImgKeyEbitender)
	opt := ebiten.DrawImageOptions{}
	screenSize := screen.Bounds().Size()
	imgSize := img.Bounds().Size()
	trans := geom.PointFFromPoint(screenSize.Sub(imgSize))
	opt.GeoM.Translate(trans.X, trans.Y)
	screen.DrawImage(img, &opt)

	drawStoryText(screen, key)
}

func drawStoryText(screen *ebiten.Image, key string) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(24, 24)
	opt.ColorScale.Scale(1, 0.8, 0.5, 1)
	drawing.DrawTextByKey(screen, key, 24, &opt)
}
