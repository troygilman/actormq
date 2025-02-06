package types

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget struct {
	game ebiten.Game
	rect image.Rectangle
}

type Container struct {
	widgets []*Widget
	rect    image.Rectangle
}

func NewContainer(widgets ...ebiten.Game) ebiten.Game {
	container := &Container{
		rect: image.Rect(0, 0, 200, 200),
	}

	for _, game := range widgets {
		container.widgets = append(container.widgets, &Widget{
			game: game,
			rect: image.Rect(0, 0, 100, 100),
		})
	}

	return container
}

func (container *Container) Update() error {
	for _, widget := range container.widgets {
		if err := widget.game.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (container *Container) Draw(screen *ebiten.Image) {
	for _, widget := range container.widgets {
		window := screen.SubImage(widget.rect).(*ebiten.Image)
		widget.game.Draw(window)
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

		widgetDx := widget.rect.Dx() + widgetDx2
		widgetDy := widget.rect.Dy() + widgetDy2

		widgetDx, widgetDy = widget.game.Layout(widgetDx, widgetDy)
	}
	return outsideWidth, outsideHeight
}
