package prologue

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
)

type Story1Drawer struct {
}

func (d *Story1Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory1)
}

type Story2Drawer struct {
	vendors []vendor
}

func (d *Story2Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory2)
	for i, v := range d.vendors {
		y := float64(i)*160 + 160
		d.drawVendor(screen, v, y)
	}
}

func (d *Story2Drawer) drawVendor(screen *ebiten.Image, v vendor, y float64) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(16, y)
	opt.ColorScale = v.ColorScale
	screen.DrawImage(drawing.Image(v.ImgKey), &opt)

	opt.ColorScale = ebiten.ColorScale{}
	opt.GeoM.Translate(136, 16)
	drawing.DrawTextByKey(screen, v.NameKey, 20, &opt)
	opt.GeoM.Translate(0, 32)
	drawing.DrawTextByKey(screen, v.DescKey, 16, &opt)
}

type vendor struct {
	NameKey    string
	ImgKey     string
	DescKey    string
	ColorScale ebiten.ColorScale
}

type Story3Drawer struct {
	managers []manager
}

func (d *Story3Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory3)
	for i, m := range d.managers {
		y := float64(i)*160 + 160
		d.drawManager(screen, m, y)
	}
}

func (d *Story3Drawer) drawManager(screen *ebiten.Image, m manager, y float64) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(16, y)
	opt.ColorScale = m.ColorScale
	screen.DrawImage(drawing.Image(m.ImgKey), &opt)

	opt.ColorScale = ebiten.ColorScale{}
	opt.GeoM.Translate(136, 8)
	drawing.DrawTextByKey(screen, m.NameKey, 20, &opt)
	opt.GeoM.Translate(0, 30)
	drawing.DrawTextByKey(screen, m.DescKey, 14, &opt)
	opt.GeoM.Translate(0, 30)
	drawing.DrawText(screen, "PERK:", 14, &opt)
	opt.GeoM.Translate(8, 18)
	drawing.DrawTextByKey(screen, m.PerkKey1, 14, &opt)
	opt.GeoM.Translate(0, 18)
	drawing.DrawTextByKey(screen, m.PerkKey2, 14, &opt)
}

type manager struct {
	NameKey    string
	ImgKey     string
	DescKey    string
	PerkKey1   string
	PerkKey2   string
	ColorScale ebiten.ColorScale
}

type Story4Drawer struct{}

func (d *Story4Drawer) Draw(screen *ebiten.Image) {
	drawStoryBG(screen, name.TextKeyStory4)
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
