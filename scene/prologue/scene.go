package prologue

import "github.com/noppikinatta/ebitenginejam03/scene"

func NewPrologueScene() scene.Scene {
	return scene.NewContainer(
		scene.NewShowImageScene(15, &Story1Drawer{}),
		scene.NewShowImageScene(15, &Story2Drawer{}),
	)
}
