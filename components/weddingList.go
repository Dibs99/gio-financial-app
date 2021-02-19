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
	bigText := material.Body1(ui.theme, "Misty & Craig\n2022-04-01")
	myInset := (100 - bigText.TextSize.V) / 2.7
	return layout.UniformInset(unit.Dp(myInset)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return bigText.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(ui.theme, "Silver").Layout(gtx)
			}),
		)
	})
}

func weddingBoxArea(gtx layout.Context) layout.Dimensions {
	const r = 10
	bounds := f32.Rect(0, 0, float32(gtx.Constraints.Min.X), float32(gtx.Constraints.Min.Y))
	defer op.Save(gtx.Ops).Load()
	clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Add(gtx.Ops)
	paint.ColorOp{Color: weddingBoxButton.currentColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func weddingBox(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	gtx.Constraints.Min.Y = gtx.Px(unit.Dp(100))
	if config.CurrentScreenSize.X > config.MaxWidth {
		gtx.Constraints.Max.X = config.MaxWidth
		gtx.Constraints.Min.X = config.MaxWidth
	}

	// size := image.Pt(text.Size.X, text.Size.Y)
	area := weddingBoxArea(gtx)
	weddingBoxButton.Layout(gtx, "weddings", area)

	return weddingText(gtx, ui)
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
