package types

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Draggable struct {
}

func (draggable *Draggable) Update() error {
	return nil
}

func (draggable Draggable) Draw(window *ebiten.Image) {
	window.Fill(color.White)
}

func (draggable Draggable) Layout(rect image.Rectangle) image.Rectangle {
	return rect
}
