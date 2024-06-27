package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
)

type ScrollText struct {
	TextKey string
	Current int
	ended   bool
}

func (s *ScrollText) Update() error {
	if s.ended {
		return nil
	}
	s.Current++
	return nil
}

func (s *ScrollText) Draw(screen *ebiten.Image) {
	const (
		fadeFrame           = 15
		scrollFrame         = 30
		stopFrame           = 60
		scrollOutStartFrame = scrollFrame + stopFrame
		endFrame            = scrollOutStartFrame + scrollFrame
		fadeoutStartFrame   = endFrame - fadeFrame
	)

	var alpha float64 = 1
	if s.Current < fadeFrame {
		alpha = float64(s.Current) / float64(fadeFrame)
	} else if s.Current > fadeoutStartFrame {
		alpha = float64(endFrame-s.Current) / float64(fadeFrame)
	}
	s.drawBG(screen, alpha)

	screenSize := geom.PointFFromPoint(screen.Bounds().Size())
	txtSize := drawing.MeasureText(lang.Text(s.TextKey), 32)
	pos := screenSize.Subtract(txtSize).Multiply(0.5)
	if s.Current < scrollFrame {
		frameDelta := float64(scrollFrame - s.Current)
		xDelta := frameDelta * frameDelta * ((screenSize.X - pos.X) / (scrollFrame * scrollFrame))
		pos.X += xDelta
	}
	if s.Current > scrollOutStartFrame {
		frameDelta := float64(s.Current - scrollOutStartFrame)
		xDelta := frameDelta * frameDelta * ((screenSize.X - pos.X) / (scrollFrame * scrollFrame))
		pos.X -= xDelta
	}
	s.drawText(screen, pos)

	if s.Current >= endFrame {
		s.ended = true
	}
}

func (s *ScrollText) drawBG(screen *ebiten.Image, alpha float64) {
	size := screen.Bounds().Size()

	a := func(v uint8) uint8 {
		return uint8(float64(v) * alpha)
	}

	c := color.RGBA{R: a(128), G: a(124), B: a(124), A: a(255)}

	vector.DrawFilledRect(screen, 0, 250, float32(size.X), float32(size.Y-500), c, false)
}

func (s *ScrollText) drawText(screen *ebiten.Image, pos geom.PointF) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(pos.X, pos.Y)

	drawing.DrawTextByKey(screen, s.TextKey, 32, &opt)
}

func (s *ScrollText) End() bool {
	return s.ended
}

func (s *ScrollText) Reset() {
	s.Current = 0
	s.ended = false
}
