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

func (tl *topicList) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
}

func (tl *topicList) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
