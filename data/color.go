package data

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/name"
)

func ColorTheme(txtKey string) ebiten.ColorScale {
	cs := ebiten.ColorScale{}

	switch txtKey {
	case name.TextKeyManager1:
		cs.Scale(0.8, 1, 1, 1)
	case name.TextKeyManager2:
		cs.Scale(1, 0.4, 0, 1)
	case name.TextKeyManager3:
		cs.Scale(1, 1, 0.25, 1)
	case name.ImgKeyVendor1:
		cs.Scale(0, 0.8, 0, 1)
	case name.ImgKeyVendor2:
		cs.Scale(0.2, 0.2, 1, 1)
	case name.ImgKeyVendor3:
		cs.Scale(1, 0.2, 0.8, 1)
	}

	return cs
}
