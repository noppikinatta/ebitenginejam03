package drawing

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	dummyImageBase = ebiten.NewImage(3, 3)

	// WhitePixel is useful to draw fill shape with DrawTriangles.
	WhitePixel = dummyImageBase.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)
