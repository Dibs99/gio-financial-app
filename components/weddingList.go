package components

import (
	"gio-test/haslett/config"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func weddingText(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	flex := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
	biggestTest := material.Body1(ui.theme, "Misty & Craig\n2022-04-01")
	return layout.UniformInset(unit.Dp((float32(gtx.Constraints.Max.Y)-biggestTest.TextSize.V)/2.7)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return biggestTest.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(ui.theme, "Silver").Layout(gtx)
			}),
		)
	})
}

func weddingBoxArea(gtx layout.Context) layout.Dimensions {
	const r = 10
	bounds := f32.Rect(0, 0, float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y))
	defer op.Push(gtx.Ops).Pop()
	clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Add(gtx.Ops)
	paint.ColorOp{Color: weddingBoxButton.currentColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func weddingBox(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	gtx.Constraints.Max.Y = 100
	gtx.Constraints.Max.X = config.MaxWidth
	gtx.Constraints.Min.X = config.MaxWidth
	// size := image.Pt(text.Size.X, text.Size.Y)
	area := weddingBoxArea(gtx)
	weddingText(gtx, ui)

	return weddingBoxButton.Layout(gtx, "weddings", area)
}

func WeddingList(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	var array []layout.Widget

	array = append(array, func(gtx layout.Context) layout.Dimensions {
		return weddingBox(gtx, ui)
	})

	return listWedding.Layout(gtx, len(array), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, array[i])
	})
}
