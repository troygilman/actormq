package types

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Draggable struct {
	rect           image.Rectangle
	dragging       bool
	cursorPosition image.Point
}

func (draggable *Draggable) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Println("PRESS", cursorX, cursorY, draggable.rect)
		if draggable.rect.Min.X <= cursorX && draggable.rect.Min.Y <= cursorY && draggable.rect.Max.X >= cursorX && draggable.rect.Max.Y >= cursorY {
			draggable.dragging = true
			draggable.cursorPosition = image.Pt(cursorX, cursorY)
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		draggable.dragging = false
	}

	if draggable.dragging {
		dx := cursorX - draggable.cursorPosition.X
		dy := cursorY - draggable.cursorPosition.Y
		draggable.rect = draggable.rect.Add(image.Pt(dx, dy))
		draggable.cursorPosition = image.Pt(cursorX, cursorY)
	}

	return nil
}

func (draggable Draggable) Draw(window *ebiten.Image) {
	window.Fill(color.White)
}

func (draggable *Draggable) Layout(rect image.Rectangle) image.Rectangle {
	if draggable.rect == image.Rect(0, 0, 0, 0) {
		draggable.rect = rect
	}
	return draggable.rect
}
