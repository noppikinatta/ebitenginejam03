package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/noppikinatta/ebitenginejam03/asset"
)

func DrawText(screen *ebiten.Image, txt string, fontSize float64, opt *ebiten.DrawImageOptions) {
	txtImg := textImage(txt, fontSize)
	screen.DrawImage(txtImg, opt)
}

func textImage(txt string, fontSize float64) *ebiten.Image {
	k := textKey{txt, fontSize}
	img, ok := textCache[k]
	if ok {
		return img
	}
	face := asset.FontFace(fontSize)
	w, h := text.Measure(txt, face, face.Metrics().HLineGap)
	img = ebiten.NewImage(int(w+1), int(h+1))
	opt := text.DrawOptions{}
	text.Draw(img, txt, face, &opt)
	textCache[k] = img
	return img
}

type textKey struct {
	Text     string
	FontSize float64
}

var textCache map[textKey]*ebiten.Image

func init() {
	textCache = make(map[textKey]*ebiten.Image)
}
