package ui

import (
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
