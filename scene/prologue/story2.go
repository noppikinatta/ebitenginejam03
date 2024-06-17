package prologue

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Story2Drawer struct{}

func (d *Story2Drawer) Draw(screen *ebiten.Image) {
	// TODO: epic story 2
	screen.Fill(color.Gray{Y: 72})
	ebitenutil.DebugPrint(screen, "story 2")
}
