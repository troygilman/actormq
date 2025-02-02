package ui

import (
	"image"

	"github.com/anthdm/hollywood/actor"
	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Update(events []any) error
	Draw(screen *ebiten.Image)
}

type WidgetContainer struct {
	widget  Widget
	rect    image.Rectangle
	engine  *actor.Engine
	client  *actor.PID
	adapter *actor.PID
	events  <-chan any
}

func (w *WidgetContainer) Update() error {
	events := pollEvents(w.events)
	return w.widget.Update(events)
}

func pollEvents(c <-chan any) []any {
	events := []any{}
	for {
		select {
		case event := <-c:
			events = append(events, event)
		default:
			return events
		}
	}
}

func (w *WidgetContainer) Draw(screen *ebiten.Image) {
	window := screen.SubImage(w.rect).(*ebiten.Image)
	w.widget.Draw(window)
}
