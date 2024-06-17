// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 noppikinatta

package scene

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	End() bool
	Reset()
}
