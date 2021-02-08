package components

import (
	"fmt"
	"image"
	"image/color"

	"gio-test/haslett/config"

	"gio-test/haslett/apiCalls"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list        = layout.List{Axis: layout.Vertical}
	listWedding = layout.List{Axis: layout.Vertical, Alignment: layout.Middle}
)

type ThisUi struct {
	theme  *material.Theme
	active int
	// list   layout.List
}

func NewUI() *ThisUi {
	ui := &ThisUi{
		theme: material.NewTheme(gofont.Collection()),
		// list: layout.List{
		//  Axis:      layout.Vertical,
		//  Alignment: layout.Middle,
		// },
	}

	return ui
}

func firstInput(gtx layout.Context, ui *ThisUi) layout.Widget {

	e := material.Editor(ui.theme, editor, "start date yyyy-MM-dd")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	return func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
		})
	}
}

func inputWidget(gtx layout.Context, ui *ThisUi) layout.Widget {

	for _, e := range lineEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			apiCalls.GetStatsWithDate(editor.Text(), e.Text)
			lineEditor.SetText("")
		}
	}

	e := material.Editor(ui.theme, lineEditor, "end date yyyy-MM-dd")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	return func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
		})
	}
}

func FinanceList(ops *op.Ops, gtx layout.Context, ui *ThisUi) layout.Dimensions {
	var (
		// uniformInset     = 16
		financeChildMaxY = 50
	)
	widgets := financeChildren(gtx, ui, financeChildMaxY)
	flex := layout.Flex{Axis: layout.Vertical}
	// return list.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
	// 	return layout.UniformInset(unit.Dp(float32(uniformInset))).Layout(gtx, widgets[i])
	// })

	return flex.Layout(gtx, widgets...)

}

func financeChildren(gtx layout.Context, ui *ThisUi, financeChildMaxY int) []layout.FlexChild {
	var array []layout.FlexChild
	array = append(array, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, firstInput(gtx, ui))
	}))
	array = append(array, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, inputWidget(gtx, ui))
	}))
	stats := apiCalls.MyStats
	for i := 0; i < len(stats.Data.ReadBankStatements); i++ {
		child := financeChild(gtx, financeChildMaxY, ui, i)
		widget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, child)
		})
		// widget := layout.Rigid(financeChild(gtx, financeChildMaxY, ui, i))
		array = append(array, widget)
	}
	return array
}

func financeChild(gtx layout.Context, financeChildMaxY int, ui *ThisUi, index int) layout.Widget {
	var stat = apiCalls.MyStats.Data.ReadBankStatements[index]
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.Y = financeChildMaxY
		gtx.Constraints.Min.Y = financeChildMaxY
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		// box
		di := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		stack := op.Push(gtx.Ops)
		clip.Rect{Max: di}.Add(gtx.Ops)
		paint.ColorOp{Color: config.Red}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Pop()

		flex3 := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}
		dims := flex3.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// circle and Category
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
				// text
				myInset := unit.Dp(float32(gtx.Constraints.Max.Y / 2))
				body := material.Body1(ui.theme, fmt.Sprintf("$%v", stat.Total))
				return layout.Inset{Right: myInset, Top: unit.Dp((float32(gtx.Constraints.Max.Y) - body.TextSize.V) / 2.5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return body.Layout(gtx)
				})
			}),
		)
		return dims
	}

}

func ErrorMessage(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	myInset := 50
	di := image.Pt(gtx.Constraints.Max.X, myInset)
	stack := op.Push(gtx.Ops)
	clip.Rect{Max: di}.Add(gtx.Ops)
	paint.ColorOp{Color: config.Red}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()
	text := material.Body1(ui.theme, apiCalls.Error)
	return layout.UniformInset(unit.Dp((float32(myInset)-text.TextSize.V)/2.5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return text.Layout(gtx)
	})
}

func FooterTabs(gtx layout.Context, screenSize int) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(float32(screenSize - FinanceButton.sizeY))}
	flex := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return FinanceButton.Layout(gtx, "finance")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return weddingButton.Layout(gtx, "wedding")
			}),
		)
	})
}
