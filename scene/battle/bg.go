package battle

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/random"
)

func CreateBG(w, h int) *ebiten.Image {
	pw := NewPixelWriter(w, h)
	rnd := rand.New(random.Source())

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := rnd.Float64()
			p := 0.001
			if r < p {
				v := byte(r * 255 / p)
				pw.Set(x, y, v, v, v, v)
			} else {
				pw.Set(x, y, 0, 0, 0, 255)
			}
		}
	}

	return pw.ToImage()
}

type PixelWriter struct {
	data   []byte
	width  int
	height int
}

func NewPixelWriter(width, height int) *PixelWriter {
	return &PixelWriter{
		data:   make([]byte, 4*width*height),
		width:  width,
		height: height,
	}
}

func (p *PixelWriter) rIdx(x, y int) int {
	return 4 * (x + p.width*y)
}

func (p *PixelWriter) Set(x, y int, r, g, b, a byte) {
	i := p.rIdx(x, y)

	p.data[i] = r
	p.data[i+1] = g
	p.data[i+2] = b
	p.data[i+3] = a
}

func (p *PixelWriter) ToImage() *ebiten.Image {
	img := ebiten.NewImage(p.width, p.height)
	img.WritePixels(p.data)
	return img
}
