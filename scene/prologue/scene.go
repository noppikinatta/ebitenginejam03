package prologue

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewPrologueScene() scene.Scene {
	s := scene.NewContainer(
		scene.NewShowImageScene(15, &Story1Drawer{}),
		scene.NewShowImageScene(15, &Story2Drawer{}),
		scene.NewShowImageScene(15, &Story3Drawer{}),
		scene.NewShowImageScene(15, &Story4Drawer{}),
	)

	s.Handlers = append(s.Handlers, &scene.ResetHandler{Key: ebiten.KeyR})

	return s
}
