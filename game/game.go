package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/name"
	negodomain "github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
	"github.com/noppikinatta/ebitenginejam03/scene/battle"
	"github.com/noppikinatta/ebitenginejam03/scene/nego"
	"github.com/noppikinatta/ebitenginejam03/scene/prologue"
)

type Game struct {
	scenes       *scene.Container
	langSwitcher *langSwitcher
}

func NewGame() *Game {
	//title := title.NewTitleScene()
	prologue := prologue.NewPrologueScene()
	negotiation, resulter := nego.NewNegotiationScene()
	battle := battle.NewBattleScene(resulter)

	scenes := scene.NewContainer(prologue, negotiation, battle)

	scenes.AddTransition(battle, prologue)

	g := Game{
		scenes:       scenes,
		langSwitcher: &langSwitcher{},
	}
	return &g
}

func (g *Game) Update() error {
	g.langSwitcher.Update()
	return g.scenes.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scenes.Draw(screen)
	g.langSwitcher.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func OrderForTest() []*negodomain.Equip {
	eqps := make([]*negodomain.Equip, 0)
	o := func(name string, improvedCount int) {
		e := &negodomain.Equip{
			Name: name, ImprovedCount: improvedCount,
		}
		eqps = append(eqps, e)
	}

	o(name.TextKeyEquip1Laser, 3)

	return eqps
}
