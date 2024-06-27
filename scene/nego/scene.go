package nego

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewNegotiationScene() (scene.Scene, func() []*nego.Equip) {
	negoScene := newNegotiationGameScene()
	s := scene.NewContainer(
		negoScene,
		scene.NewShowImageScene(15, newSpecDrawer(negoScene.Result)),
	)

	s.Handlers = append(s.Handlers, &scene.ResetHandler{Key: ebiten.KeyR})

	return s, negoScene.Result
}
