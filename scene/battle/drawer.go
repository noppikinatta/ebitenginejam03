package battle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/shooter"
)

type visibleEntityDrawer interface {
	Draw(screen *ebiten.Image, entity shooter.VisibleEntity)
}

type laserDrawer struct{}

func (d *laserDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}

type missileDrawer struct{}

func (d *missileDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}

type harakiriDrawer struct{}

func (d *harakiriDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}

type barrierDrawer struct{}

func (d *barrierDrawer) Draw(screen *ebiten.Image, entity shooter.VisibleEntity) {

}
