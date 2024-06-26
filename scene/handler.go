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
