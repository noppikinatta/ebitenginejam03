// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 noppikinatta

package scene

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var faderGlobalCacheImage *ebiten.Image

type FadeIn struct {
	fade
}

func NewFadeIn(frames int) *FadeIn {
	f := FadeIn{
		fade: fade{
			frames: frames,
		},
	}

	return &f
}

func (f *FadeIn) Draw(target *ebiten.Image) {
	f.fade.Draw(target, f.alpha())
}

func (f *FadeIn) alpha() float32 {
	if f.frames == 0 {
		return 0
	}

	a := 1.0 - float32(f.current)/float32(f.frames)
	if a < 0 {
		a = 0
	}
	return a
}

type FadeOut struct {
	fade
}

func NewFadeOut(frames int) *FadeOut {
	f := FadeOut{
		fade: fade{
			frames: frames,
		},
	}

	return &f
}

func (f *FadeOut) Draw(screen *ebiten.Image) {
	f.fade.Draw(screen, f.alpha())
}

func (f *FadeOut) alpha() float32 {
	if f.frames == 0 {
		return 0
	}

	a := float32(f.current) / float32(f.frames)
	if a > 1 {
		a = 1
	}
	return a
}

type fade struct {
	current int
	frames  int
}

func (f *fade) Reset() {
	f.current = 0
}

func (f *fade) Update() error {
	if !f.End() {
		f.current++
	}
	return nil
}

func (f *fade) Draw(screen *ebiten.Image, alpha float32) {
	if alpha == 0 {
		return
	}
	c := f.cache(screen.Bounds().Size())

	opt := ebiten.DrawImageOptions{}
	opt.ColorScale.SetA(alpha)

	screen.DrawImage(c, &opt)
}

func (f *fade) cache(size image.Point) *ebiten.Image {
	if f.shouldUpdateCache(size) {
		f.updateCache(size)
	}

	return faderGlobalCacheImage
}

func (f *fade) shouldUpdateCache(size image.Point) bool {
	if faderGlobalCacheImage == nil {
		return true
	}
	cacheSize := faderGlobalCacheImage.Bounds().Size()
	return size == cacheSize
}

func (f *fade) updateCache(size image.Point) {
	if faderGlobalCacheImage != nil {
		faderGlobalCacheImage.Deallocate()
	}
	faderGlobalCacheImage = ebiten.NewImage(size.X, size.Y)
	faderGlobalCacheImage.Fill(color.Black)
}

func (f *fade) End() bool {
	if f.frames == 0 {
		return true
	}
	return f.current > f.frames
}
