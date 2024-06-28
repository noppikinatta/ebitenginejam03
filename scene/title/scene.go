package title

import "github.com/noppikinatta/ebitenginejam03/scene"

func NewTitleScene() scene.Scene {
	return scene.NewShowImageScene(24, &TitleDrawer{})
}
