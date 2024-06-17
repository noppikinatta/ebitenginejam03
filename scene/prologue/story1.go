package prologue

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Story1Drawer struct {
}

func (d *Story1Drawer) Draw(screen *ebiten.Image) {
	// TODO: epic story 1
	screen.Fill(color.Gray{Y: 96})
	ebitenutil.DebugPrint(screen, "story 1")
}
