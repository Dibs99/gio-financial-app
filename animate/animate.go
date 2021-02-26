package animate

import (
	"gio-test/haslett/components"
	"gio-test/haslett/config"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type FadeInStr struct {
	callback  func()
	startTime time.Time
	duration  int
}

var duration = 1 * time.Second

func FadeInValue(gtx layout.Context, value color.NRGBA, valueToFadeTo uint8) {

	elapsed := time.Now().Sub(components.TimeTrigger)
	progress := elapsed.Seconds() / duration.Seconds()

	if progress < 1 {
		// The progress bar hasn’t yet finished animating.
		op.InvalidateOp{}.Add(gtx.Ops)
	} else {
		progress = 1
	}
	myStack := op.Save(gtx.Ops)
	s := config.Background
	d := 255 - (progress * 254)
	s.A = uint8(d)
	paint.Fill(gtx.Ops, s)
	myStack.Load()

}

// type callback func(layout.Dimensions) layout.Dimensions

// func FadeInContainer(gtx layout.Context, callback callback, w layout.Dimensions) {

// 	elapsed := time.Now().Sub(startTime)
// 	progress := elapsed.Seconds() / duration.Seconds()
// 	if progress < 1 {
// 		// The progress bar hasn’t yet finished animating.
// 		op.InvalidateOp{}.Add(gtx.Ops)
// 	} else {
// 		progress = 1
// 	}
// 	myStack := op.Save(gtx.Ops)
// 	callback(w)
// 	paint.Fill(gtx.Ops, color.NRGBA{R: 187, G: 169, B: 169, A: uint8(progress * 254)})
// 	myStack.Load()

// }
