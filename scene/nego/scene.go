package nego

import "github.com/noppikinatta/ebitenginejam03/scene"

func NewNegotiationScene() scene.Scene {
	return scene.NewContainer(
		scene.NewFadeIn(15),
		NewNegotiationGameScene(),
		scene.NewFadeOut(15),
	)
}
