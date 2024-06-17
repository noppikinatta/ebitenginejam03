package title

import "github.com/noppikinatta/ebitenginejam03/scene"

func NewTitleScene() scene.Scene {
	return scene.NewShowImageScene(15, &TitleDrawer{})
}
