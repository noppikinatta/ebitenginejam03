package font

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed Mplus2-Regular.ttf
	mplus2regularttf []byte

	Mplus2RegularFaceSource *text.GoTextFaceSource

	fontFaces map[float64]text.Face
)

func init() {
	var err error

	Mplus2RegularFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(mplus2regularttf))
	if err != nil {
		log.Fatal(err)
	}

	fontFaces = make(map[float64]text.Face)
}

func FontFace(size float64) text.Face {
	f, ok := fontFaces[size]
	if !ok {
		f = &text.GoTextFace{
			Source: Mplus2RegularFaceSource,
			Size:   size,
		}
		fontFaces[size] = f
	}

	return f
}
