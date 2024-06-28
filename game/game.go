package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/name"
	negodomain "github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/scene"
	"github.com/noppikinatta/ebitenginejam03/scene/battle"
	"github.com/noppikinatta/ebitenginejam03/scene/nego"
	"github.com/noppikinatta/ebitenginejam03/scene/prologue"
	"github.com/noppikinatta/ebitenginejam03/scene/title"
)

type Game struct {
	scenes       *scene.Container
	langSwitcher *langSwitcher
}

func NewGame() *Game {
	title := title.NewTitleScene()
	prologueScene := prologue.NewPrologueScene()
	negotiationScene, resulter := nego.NewNegotiationScene()
	battleScene := battle.NewBattleScene(resulter)

	scenes := scene.NewContainer(title, prologueScene, negotiationScene, battleScene)
	scenes.Handlers = append(scenes.Handlers, &scene.LongPressResetHandler{Key: ebiten.KeyR, WaitUntilReset: 120})

	scenes.AddTransition(battleScene, prologueScene)

	// debug code
	//battleScene = battle.NewBattleScene(OrderForTest)
	//scenes = scene.NewContainer(battleScene)

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

	improvedCount := 2

	o(name.TextKeyEquip1Laser, improvedCount)
	o(name.TextKeyEquip2Missile, improvedCount)
	o(name.TextKeyEquip3Harakiri, improvedCount)
	o(name.TextKeyEquip4Barrier, improvedCount)
	o(name.TextKeyEquip5Armor, improvedCount)
	o(name.TextKeyEquip6Exhaust, improvedCount)
	o(name.TextKeyEquip7Stonehenge, improvedCount)
	o(name.TextKeyEquip8Sushibar, improvedCount)
	o(name.TextKeyEquip9Operahouse, improvedCount)

	return eqps
}
