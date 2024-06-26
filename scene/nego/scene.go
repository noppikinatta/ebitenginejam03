package nego

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewNegotiationScene() (scene.Scene, func() []*nego.Equip) {
	negoScene := newNegotiationGameScene()
	s := scene.NewContainer(
		scene.NewFadeIn(15),
		negoScene,
		scene.NewFadeOut(15),
	)

	s.Handlers = append(s.Handlers, &scene.ResetHandler{Key: ebiten.KeyR})

	return s, negoScene.Result
}
