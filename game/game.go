package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/scene"
	"github.com/noppikinatta/ebitenginejam03/scene/nego"
	"github.com/noppikinatta/ebitenginejam03/scene/prologue"
	"github.com/noppikinatta/ebitenginejam03/scene/title"
)

type Game struct {
	scenes *scene.Container
}

func NewGame() *Game {
	title := title.NewTitleScene()
	prologue := prologue.NewPrologueScene()
	negotiation := nego.NewNegotiationScene()

	scenes := scene.NewContainer(title, prologue, negotiation)

	scenes.AddTransition(negotiation, title)

	g := Game{
		scenes: scenes,
	}
	return &g
}

func (g *Game) Update() error {
	return g.scenes.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scenes.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}
