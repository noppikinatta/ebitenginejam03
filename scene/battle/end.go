package battle

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
	"github.com/noppikinatta/ebitenginejam03/name"
)

type battleEndScene struct {
	resultFn func() bool
	frames   int
	current  int
}

func (s *battleEndScene) Update() error {
	if s.End() {
		return nil
	}
	s.current++
	return nil
}

func (s *battleEndScene) Draw(screen *ebiten.Image) {
	var titleKey string
	var epilogueKey string

	if s.resultFn() {
		titleKey = name.TextKeyShooterResult1
		epilogueKey = name.TextKeyShooterResult1Epilogue
	} else {
		titleKey = name.TextKeyShooterResult2
		epilogueKey = name.TextKeyShooterResult2Epilogue
	}

	alpha := float64(s.current) / float64(s.frames)
	bgAlpha := alpha * 0.5

	size := screen.Bounds().Size()
	clr := color.RGBA{A: uint8(255 * bgAlpha)}
	vector.DrawFilledRect(screen, 0, 0, float32(size.X), float32(size.Y), clr, false)

	s.drawTitle(screen, titleKey, alpha)
	s.drawEpilogue(screen, epilogueKey, alpha)
}

func (s *battleEndScene) drawTitle(screen *ebiten.Image, titleKey string, alpha float64) {
	screenSize := geom.PointFFromPoint(screen.Bounds().Size())

	var fontSize float64 = 64
	titleSize := drawing.MeasureText(lang.Text(titleKey), fontSize)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate((screenSize.X-titleSize.X)*0.5, 100)
	opt.ColorScale.ScaleAlpha(float32(alpha))
	drawing.DrawTextByKey(screen, titleKey, fontSize, &opt)
}

func (s *battleEndScene) drawEpilogue(screen *ebiten.Image, epilogueKey string, alpha float64) {
	screenSize := geom.PointFFromPoint(screen.Bounds().Size())

	var fontSize float64 = 32
	titleSize := drawing.MeasureText(lang.Text(epilogueKey), fontSize)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate((screenSize.X-titleSize.X)*0.5, 280)
	opt.ColorScale.ScaleAlpha(float32(alpha))
	drawing.DrawTextByKey(screen, epilogueKey, fontSize, &opt)
}

func (s *battleEndScene) End() bool {
	return s.current >= s.frames
}

func (s *battleEndScene) Reset() {
	s.current = 0
}
