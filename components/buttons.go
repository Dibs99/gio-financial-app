package components

import (
	"gio-test/haslett/config"
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

var (
	button        = new(widget.Clickable)
	buttonState   = true
	financeButton = &Button{pressed: false, currentColor: config.Blue, initialColor: config.Blue, hoverColor: config.Red, sizeX: 100, sizeY: 50}
	weddingButton = &Button{pressed: false, currentColor: config.Green, initialColor: config.Green, hoverColor: config.Red, sizeX: 100, sizeY: 50}
)

type Button struct {
	pressed      bool
	currentColor color.NRGBA
	hoverColor   color.NRGBA
	initialColor color.NRGBA
	sizeX        int
	sizeY        int
}

func (b *Button) Layout(gtx layout.Context, screen string) layout.Dimensions {
	// Avoid affecting the input tree with pointer events.
	defer op.Push(gtx.Ops).Pop()

	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				config.CurrentScreen = screen
			case pointer.Enter:
				b.currentColor = b.hoverColor
			case pointer.Leave:
				b.currentColor = b.initialColor
			}
		}
	}

	// Confine the area for pointer events.
	pointer.Rect(image.Rect(0, 0, b.sizeX, b.sizeY)).Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)

	// Draw the button.
	return ColorBox(gtx, image.Pt(b.sizeX, b.sizeY), b.currentColor)
}

func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
