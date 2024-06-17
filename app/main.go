package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/game"
)

func main() {
	ebiten.SetWindowSize(800, 640)
	ebiten.SetWindowTitle("Ebiten Game Jam")

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
