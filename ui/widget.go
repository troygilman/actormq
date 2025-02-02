package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type WidgetContainer struct {
	widget Widget
	rect   image.Rectangle
}

func NewWidgetContainer(widget Widget, rect image.Rectangle) *WidgetContainer {
	return &WidgetContainer{
		widget: widget,
		rect:   rect,
	}
}

func (w *WidgetContainer) Draw(screen *ebiten.Image) {
	window := screen.SubImage(w.rect).(*ebiten.Image)
	w.widget.Draw(window)
}
