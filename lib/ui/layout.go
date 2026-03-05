package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layout struct {
	Drawer
	Translate image.Point
}

func (l *Layout) Draw(screen *ebiten.Image) {
	l.Drawer.Draw(screen, l.Translate)
}
