package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Handler interface {
	Handle(scene Scene)
}

type ResetHandler struct {
	Key  ebiten.Key
	keys []ebiten.Key
}

func (h *ResetHandler) Handle(scene Scene) {
	h.keys = h.keys[:0]
	h.keys = inpututil.AppendJustPressedKeys(h.keys)
	for _, k := range h.keys {
		if k == h.Key {
			scene.Reset()
			return
		}
	}
}

type LongPressResetHandler struct {
	Key            ebiten.Key
	keys           []ebiten.Key
	count          int
	WaitUntilReset int
}

func (h *LongPressResetHandler) Handle(scene Scene) {
	if h.WaitUntilReset == 0 {
		return
	}

	if h.keyPressed() {
		h.count++
	} else {
		h.count = 0
	}

	if h.count < h.WaitUntilReset {
		return
	}

	scene.Reset()
}

func (h *LongPressResetHandler) keyPressed() bool {
	h.keys = h.keys[:0]
	h.keys = inpututil.AppendPressedKeys(h.keys)
	for _, k := range h.keys {
		if k == h.Key {
			return true
		}
	}

	return false
}
