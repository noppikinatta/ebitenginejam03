// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 noppikinatta

package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
	scenes      []Scene
	currentIdx  int
	transitions map[int]int
}

func NewContainer(scenes ...Scene) *Container {
	c := Container{
		scenes:      scenes,
		currentIdx:  0,
		transitions: make(map[int]int),
	}
	return &c
}

func (c *Container) AddTransition(from, to Scene) error {
	fi, ok := c.indexOf(from)
	if !ok {
		return errors.New("from scene for transition is not in container")
	}

	ti, ok := c.indexOf(to)
	if !ok {
		return errors.New("to scene for transition is not in container")
	}

	c.transitions[fi] = ti
	return nil
}

func (c *Container) indexOf(s Scene) (int, bool) {
	for i := range c.scenes {
		if s == c.scenes[i] {
			return i, true
		}
	}

	return 0, false
}

func (c *Container) Update() error {
	if c.End() {
		return nil
	}

	current, ok := c.current()
	if !ok {
		return nil
	}

	return current.Update()
}

func (c *Container) Draw(screen *ebiten.Image) {
	current, ok := c.current()
	if !ok {
		return
	}
	current.Draw(screen)
	if current.End() {
		c.next()
	}
}

func (c *Container) End() bool {
	return c.currentIdx >= len(c.scenes)
}

func (c *Container) Reset() {
	for _, s := range c.scenes {
		s.Reset()
	}
	c.currentIdx = 0
}

func (c *Container) current() (Scene, bool) {
	if c.currentIdx < 0 || c.currentIdx >= len(c.scenes) {
		return nil, false
	}

	return c.scenes[c.currentIdx], true
}

func (c *Container) next() {
	c.currentIdx = c.getNext()
	s, ok := c.current()
	if ok {
		s.Reset()
	}
}

func (c *Container) getNext() int {
	idx, ok := c.transitions[c.currentIdx]
	if ok {
		return idx
	}

	idx = c.currentIdx + 1
	return idx
}
