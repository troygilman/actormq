package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget struct {
	game   ebiten.Game
	rect   image.Rectangle
	window *ebiten.Image
}

func NewWidget(game ebiten.Game, rect image.Rectangle) *Widget {
	return &Widget{
		game: game,
		rect: rect,
	}
}

func (w *Widget) Draw(screen *ebiten.Image) {
	if w.window == nil {
		w.window = screen.SubImage(w.rect).(*ebiten.Image)
	}
	w.game.Draw(w.window)
}
