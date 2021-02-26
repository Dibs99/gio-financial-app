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
	listWedding = new(layout.List)
	// test        = &listWedding
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
		financeChildMaxY float32 = 50
	)
	widgets := financeChildren(gtx, ui, financeChildMaxY)
	flex := layout.Flex{Axis: layout.Vertical}

	return flex.Layout(gtx, widgets...)

}

func financeChildren(gtx layout.Context, ui *ThisUi, financeChildMaxY float32) []layout.FlexChild {
	var array []layout.FlexChild
	array = append(array, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, firstInput(gtx, ui))
	}))
	array = append(array, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, inputWidget(gtx, ui))
	}))
	array = append(array, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(ui.theme, "*defaults to the last 6 months of records").Layout(gtx)
		})
	}))
	stats := apiCalls.MyStats
	for i := 0; i < len(stats.Data.ReadBankStatements); i++ {
		child := financeChild(gtx, financeChildMaxY, ui, i)
		widget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, child)
		})
		array = append(array, widget)
	}
	return array
}

func financeChild(gtx layout.Context, financeChildMaxY float32, ui *ThisUi, index int) layout.Widget {
	var stat = apiCalls.MyStats.Data.ReadBankStatements[index]
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.Y = gtx.Px(unit.Dp(financeChildMaxY))
		gtx.Constraints.Min.Y = gtx.Px(unit.Dp(financeChildMaxY))
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		// box
		di := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		stack := op.Save(gtx.Ops)
		clip.Rect{Max: di}.Add(gtx.Ops)
		paint.ColorOp{Color: config.Red}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()

		flex3 := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}
		dims := flex3.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// circle and Category
				max := gtx.Constraints.Max.Y / 2
				myInset := unit.Dp(financeChildMaxY)
				body := material.Body1(ui.theme, stat.MyCategory)
				return layout.Inset{Left: myInset, Top: unit.Dp((financeChildMaxY - body.TextSize.V) / 2.5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					roundness := float32(max / 2)
					stack := op.Save(gtx.Ops)
					op.Offset(f32.Pt(float32(-max-5), 0)).Add(gtx.Ops)
					clip.RRect{Rect: f32.Rect(0, 0, float32(max), float32(max)), SE: roundness, SW: roundness, NW: roundness, NE: roundness}.Add(gtx.Ops)
					paint.ColorOp{Color: stat.Colour}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					stack.Load()
					return body.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// text
				text := fmt.Sprintf("$%v", stat.Total)
				if PieChartAreaButton.pressed {
					text = fmt.Sprintf("%v%%", stat.Percentage)
				}
				myInset := unit.Dp(financeChildMaxY / 2)
				body := material.Body1(ui.theme, text)
				return layout.Inset{Right: myInset, Top: unit.Dp((financeChildMaxY - body.TextSize.V) / 2.5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return body.Layout(gtx)
				})
			}),
		)
		return dims
	}

}

func ErrorMessage(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	var myInset float32 = 50
	di := image.Pt(gtx.Constraints.Max.X, gtx.Px(unit.Dp(myInset)))
	stack := op.Save(gtx.Ops)
	clip.Rect{Max: di}.Add(gtx.Ops)
	paint.ColorOp{Color: config.Red}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Load()
	text := material.Body1(ui.theme, apiCalls.Error)
	text.MaxLines = 1
	return layout.UniformInset(unit.Dp((float32(myInset)-text.TextSize.V)/2.5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return text.Layout(gtx)
	})
}

func FooterTabs(gtx layout.Context, screenSize image.Point, ui *ThisUi) layout.Dimensions {
	// inset :=
	flex := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}

	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return FinanceButton.Layout(gtx, "finance")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return weddingButton.Layout(gtx, "weddings")
		}),
	)

}
