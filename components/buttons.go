package components

import (
	"gio-test/haslett/apiCalls"
	"gio-test/haslett/config"
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

var (
	button                   = new(widget.Clickable)
	buttonState              = true
	FinanceButton            = Button{pressed: false, currentColor: config.Blue, initialColor: config.Blue, hoverColor: config.Red, sizeX: 100, sizeY: 50}
	weddingButton            = Button{pressed: false, currentColor: config.Green, initialColor: config.Green, hoverColor: config.Red, sizeX: 100, sizeY: ButtonSizeY, callBack: apiCalls.GetAllWeddings}
	weddingBoxButton         = AreaButton{pressed: false, currentColor: config.Red, initialColor: config.Red, hoverColor: config.OffWhite}
	ButtonSizeY      float32 = 50
	// weddingBoxButton =
)

type Button struct {
	pressed      bool
	currentColor color.NRGBA
	hoverColor   color.NRGBA
	initialColor color.NRGBA
	sizeX        float32
	sizeY        float32
	callBack     func()
}

func (b *Button) Layout(gtx layout.Context, screen string) layout.Dimensions {
	// Avoid affecting the input tree with pointer events.
	defer op.Save(gtx.Ops).Load()

	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				config.CurrentScreen = screen
				// b.callBack()
			case pointer.Enter:
				b.currentColor = b.hoverColor
			case pointer.Leave:
				b.currentColor = b.initialColor
			}
		}
	}
	sizeX := gtx.Px(unit.Dp(b.sizeX))
	sizeY := gtx.Px(unit.Dp(b.sizeY))
	// Confine the area for pointer events.
	pointer.Rect(image.Rect(0, 0, sizeX, sizeY)).Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)

	// Draw the button.
	return ColorBox(gtx, image.Pt(sizeX, sizeY), b.currentColor)
}

func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	dr := image.Rectangle{Max: size}
	defer op.Save(gtx.Ops).Load()
	clip.Rect(dr).Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

type AreaButton struct {
	pressed      bool
	currentColor color.NRGBA
	hoverColor   color.NRGBA
	initialColor color.NRGBA
}

func (b *AreaButton) Layout(gtx layout.Context, screen string, area layout.Dimensions) layout.Dimensions {
	// Avoid affecting the input tree with pointer events.
	defer op.Save(gtx.Ops).Load()

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
	pointer.Rect(image.Rect(0, 0, area.Size.X, area.Size.Y)).Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)

	return area
}
