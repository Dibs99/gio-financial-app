package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func financeList(ops *op.Ops, gtx layout.Context, ui *UI) layout.Dimensions {
	var (
		uniformInset     = 16
		financeChildMaxY = 50
	)

	widgets := financeChildren(gtx, ui, financeChildMaxY)
	return list.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(float32(uniformInset))).Layout(gtx, widgets[i])
	})

}

func financeChildren(gtx layout.Context, ui *UI, financeChildMaxY int) []layout.Widget {
	var array []layout.Widget
	stats := myStats
	for i := 0; i < len(stats.Data.ReadBankStatements); i++ {

		widget := financeChild(gtx, financeChildMaxY, ui, i)
		array = append(array, widget)
	}
	return array
}

func financeChild(gtx layout.Context, financeChildMaxY int, ui *UI, index int) layout.Widget {
	var stat = myStats.Data.ReadBankStatements[index]
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.Y = financeChildMaxY
		gtx.Constraints.Min.Y = financeChildMaxY
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		di := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		stack := op.Push(gtx.Ops)
		clip.Rect{Max: di}.Add(gtx.Ops)
		paint.ColorOp{Color: red}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Pop()

		flex3 := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
		dims := flex3.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				max := gtx.Constraints.Max.Y / 2
				myInset := unit.Dp(float32(gtx.Constraints.Max.Y))
				body := material.Body1(ui.theme, stat.Category)
				return layout.Inset{Left: myInset, Top: unit.Dp((float32(gtx.Constraints.Max.Y) - body.TextSize.V) / 2.5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					roundness := float32(max / 2)
					stack := op.Push(gtx.Ops)
					op.Offset(f32.Pt(float32(-max-5), 0)).Add(gtx.Ops)
					clip.RRect{Rect: f32.Rect(0, 0, float32(max), float32(max)), SE: roundness, SW: roundness, NW: roundness, NE: roundness}.Add(gtx.Ops)
					paint.ColorOp{Color: stat.Colour.(color.NRGBA)}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					stack.Pop()
					return body.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				myInset := unit.Dp(float32(gtx.Constraints.Max.Y / 2))
				body := material.Body1(ui.theme, fmt.Sprintf("$%v", stat.Total))
				return layout.Inset{Right: myInset, Top: unit.Dp((float32(gtx.Constraints.Max.Y) - body.TextSize.V) / 2.5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return body.Layout(gtx)
				})
			}),
		)

		// di2 := dims.Size
		// stack2 := op.Push(gtx.Ops)
		// clip.Rect{Max: di2}.Add(gtx.Ops)
		// paint.ColorOp{Color: green}.Add(gtx.Ops)
		// paint.PaintOp{}.Add(gtx.Ops)
		// stack2.Pop()
		return dims
	}

}

func footerTabs(gtx layout.Context, screenSize int) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(float32(screenSize - financeButton.sizeY))}
	flex := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return financeButton.Layout(gtx, "finance")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return weddingButton.Layout(gtx, "wedding")
			}),
		)
	})
}
