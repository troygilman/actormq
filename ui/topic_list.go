package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type topicList struct {
}

func (tl *topicList) Update() error {
	return nil
}

func (tl *topicList) Draw(window *ebiten.Image) {
	window.Fill(color.White)
}

func (tl *topicList) Layout(rect image.Rectangle) image.Rectangle {
	return rect
}
