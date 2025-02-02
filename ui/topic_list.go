package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type topicList struct {
}

func (tl *topicList) Update(events []any) error {
	return nil
}

func (tl *topicList) Draw(window *ebiten.Image) {
	window.Fill(color.White)
}
