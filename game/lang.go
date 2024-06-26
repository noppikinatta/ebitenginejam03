package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
)

type langSwitcher struct {
	Alpha       float32
	Keys        []ebiten.Key
	CurrentLang string
}

func (s *langSwitcher) Update() {
	s.Keys = s.Keys[:0]
	s.Keys = inpututil.AppendJustPressedKeys(s.Keys)
	for _, k := range s.Keys {
		if k == ebiten.KeyL {
			s.CurrentLang = lang.Switch()
			s.Alpha = 1
			break
		}
	}

	s.Alpha -= 0.01
	if s.Alpha < 0 {
		s.Alpha = 0
	}
}

func (s *langSwitcher) Draw(screen *ebiten.Image) {
	if s.Alpha == 0 {
		return
	}

	idxs := []uint16{0, 1, 2, 0, 2, 3}

	topLeft := geom.PointF{
		X: 0,
		Y: 590,
	}

	bottomRight := geom.PointF{
		X: 200,
		Y: 620,
	}

	v := ebiten.Vertex{
		ColorA: s.Alpha,
	}
	vertices := make([]ebiten.Vertex, 4)
	v.DstX = float32(topLeft.X)
	v.DstY = float32(topLeft.Y)
	vertices[0] = v
	v.DstX = float32(topLeft.X)
	v.DstY = float32(bottomRight.Y)
	vertices[1] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(bottomRight.Y)
	v.ColorA = 0
	vertices[2] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(topLeft.Y)
	v.ColorA = 0
	vertices[3] = v

	ropt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &ropt)

	topt := ebiten.DrawImageOptions{}
	topt.GeoM.Translate(0, 592)
	topt.ColorScale.Scale(s.Alpha, s.Alpha, s.Alpha, s.Alpha)
	drawing.DrawText(screen, fmt.Sprintf("Current Language: %s", s.CurrentLang), 16, &topt)
}
