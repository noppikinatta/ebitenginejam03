package nego

import (
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewNegotiationScene() (scene.Scene, func() []*nego.Equip) {
	negoScene := newNegotiationGameScene()
	return scene.NewContainer(
		scene.NewFadeIn(15),
		negoScene,
		scene.NewFadeOut(15),
	), negoScene.Result
}
