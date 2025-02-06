package types

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Update() error
	Draw(*ebiten.Image)
	Layout(image.Rectangle) image.Rectangle
}

func NewWidget(widget Widget, rect image.Rectangle) *WidgetContainer {
	return &WidgetContainer{
		widget: widget,
		rect:   rect,
	}
}

type WidgetContainer struct {
	widget Widget
	rect   image.Rectangle
}

type Container struct {
	widgets []*WidgetContainer
	rect    image.Rectangle
}

func NewContainer(widgets ...*WidgetContainer) ebiten.Game {
	container := &Container{
		rect:    image.Rect(0, 0, 200, 200),
		widgets: widgets,
	}
	return container
}

func (container *Container) Update() error {
	for _, widget := range container.widgets {
		if err := widget.widget.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (container *Container) Draw(screen *ebiten.Image) {
	for _, widget := range container.widgets {
		window := screen.SubImage(widget.rect).(*ebiten.Image)
		widget.widget.Draw(window)
	}
}

func (container *Container) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	dx2 := outsideWidth - container.rect.Dx()
	dy2 := outsideHeight - container.rect.Dy()

	dx3 := dx2 / outsideWidth
	dy3 := dy2 / outsideHeight

	container.rect = image.Rect(0, 0, outsideWidth, outsideHeight)

	for _, widget := range container.widgets {
		widgetDx2 := dx3 * widget.rect.Dx()
		widgetDy2 := dy3 * widget.rect.Dy()
		rect := widget.rect.Add(image.Pt(widgetDx2, widgetDy2))

		widget.rect = widget.widget.Layout(rect)
	}
	return outsideWidth, outsideHeight
}
