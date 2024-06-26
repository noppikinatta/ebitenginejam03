package battle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewBattleScene(orderer func() []*nego.Equip) scene.Scene {
	s := scene.NewContainer(
		scene.NewFadeIn(15),
		newBattleGameScene(orderer),
		scene.NewFadeOut(15),
	)

	s.Handlers = append(s.Handlers, &scene.ResetHandler{Key: ebiten.KeyR})

	return s
}
