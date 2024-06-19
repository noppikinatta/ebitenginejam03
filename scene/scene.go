// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 noppikinatta

package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	End() bool
	Reset()
}

type Drawer interface {
	Draw(screen *ebiten.Image)
}

type showImageScene struct {
	fadeScenes *Container
	drawer     Drawer
}

func NewShowImageScene(frames int, drawer Drawer) Scene {
	fadeScenes := NewContainer(
		NewFadeIn(frames),
		&WaitClick{},
		NewFadeOut(frames),
	)

	return &showImageScene{
		fadeScenes: fadeScenes,
		drawer:     drawer,
	}
}

func (s *showImageScene) Update() error {
	return s.fadeScenes.Update()
}

func (s *showImageScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen)
	s.fadeScenes.Draw(screen)
}

func (s *showImageScene) End() bool {
	return s.fadeScenes.End()
}

func (s *showImageScene) Reset() {
	s.fadeScenes.Reset()
}
