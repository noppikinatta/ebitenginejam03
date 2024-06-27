package prologue

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/data"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

func NewPrologueScene() scene.Scene {
	s := scene.NewContainer(
		scene.NewShowImageScene(15, &Story1Drawer{}),
		scene.NewShowImageScene(15, &Story2Drawer{
			vendors: []vendor{
				{
					NameKey:    name.TextKeyVendor1,
					ImgKey:     name.ImgKeyVendor1,
					DescKey:    name.TextKeyVendor1Desc,
					ColorScale: data.ColorTheme(name.TextKeyVendor1),
				},
				{
					NameKey:    name.TextKeyVendor2,
					ImgKey:     name.ImgKeyVendor2,
					DescKey:    name.TextKeyVendor2Desc,
					ColorScale: data.ColorTheme(name.TextKeyVendor2),
				},
				{
					NameKey:    name.TextKeyVendor3,
					ImgKey:     name.ImgKeyVendor3,
					DescKey:    name.TextKeyVendor3Desc,
					ColorScale: data.ColorTheme(name.TextKeyVendor3),
				},
			},
		}),
		scene.NewShowImageScene(15, &Story3Drawer{
			managers: []manager{
				{
					NameKey:    name.TextKeyManager1,
					ImgKey:     name.ImgKeyManager1,
					DescKey:    name.TextKeyManager1Desc,
					PerkKey1:   name.TextKeyManager1Perk1,
					PerkKey2:   name.TextKeyManager1Perk2,
					ColorScale: data.ColorTheme(name.TextKeyManager1),
				},
				{
					NameKey:    name.TextKeyManager2,
					ImgKey:     name.ImgKeyManager2,
					DescKey:    name.TextKeyManager2Desc,
					PerkKey1:   name.TextKeyManager2Perk1,
					PerkKey2:   name.TextKeyManager2Perk2,
					ColorScale: data.ColorTheme(name.TextKeyManager2),
				},
				{
					NameKey:    name.TextKeyManager3,
					ImgKey:     name.ImgKeyManager3,
					DescKey:    name.TextKeyManager3Desc,
					PerkKey1:   name.TextKeyManager3Perk1,
					PerkKey2:   name.TextKeyManager3Perk2,
					ColorScale: data.ColorTheme(name.TextKeyManager3),
				},
			},
		}),
		scene.NewShowImageScene(15, &Story4Drawer{}),
	)

	s.Handlers = append(s.Handlers, &scene.ResetHandler{Key: ebiten.KeyR})

	return s
}
