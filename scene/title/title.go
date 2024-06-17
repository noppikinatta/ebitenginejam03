package title

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TitleDrawer struct {
}

func (d *TitleDrawer) Draw(screen *ebiten.Image) {
	// TODO: nice title image
	screen.Fill(color.Gray{Y: 128})
	ebitenutil.DebugPrint(screen, "title")
}
