package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type WaitClick struct {
	clicked bool
}

func (w *WaitClick) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		w.clicked = true
	}
	return nil
}

func (w *WaitClick) Draw(screen *ebiten.Image) {
	// do nothing
}

func (w *WaitClick) End() bool {
	return w.clicked
}

func (w *WaitClick) Reset() {
	w.clicked = false
}
