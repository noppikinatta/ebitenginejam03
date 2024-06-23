package battle

import (
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewBattleScene(orderer func() []*nego.Equip) scene.Scene {
	return scene.NewContainer(
		scene.NewFadeIn(15),
		newBattleGameScene(orders),
		scene.NewFadeOut(15),
	)
}
