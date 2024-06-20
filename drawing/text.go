package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/noppikinatta/ebitenginejam03/asset/font"
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
	face := font.FontFace(fontSize)
	w, h := text.Measure(txt, face, lineSpacing(fontSize))
	img = ebiten.NewImage(int(w+1), int(h+1))
	opt := text.DrawOptions{}
	text.Draw(img, txt, face, &opt)
	textCache[k] = img
	return img
}

func lineSpacing(fontSize float64) float64 {
	spacing, ok := lineSpacingCache[fontSize]
	if ok {
		return spacing
	}

	face := font.FontFace(fontSize)
	_, h := text.Measure("A", face, 0)
	spacing = h * 1.5
	lineSpacingCache[fontSize] = spacing
	return spacing
}

type textKey struct {
	Text     string
	FontSize float64
}

var textCache map[textKey]*ebiten.Image

var lineSpacingCache map[float64]float64

func init() {
	textCache = make(map[textKey]*ebiten.Image)
	lineSpacingCache = make(map[float64]float64)
}
