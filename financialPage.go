package main

import (
	"gio-test/haslett/apiCalls"
	"gio-test/haslett/components"

	"gioui.org/layout"
)

func financeLayout(axis layout.Axis, gtx layout.Context, stats apiCalls.NewStats, Ui *components.ThisUi) layout.Dimensions {
	flex := layout.Flex{Axis: axis, Alignment: layout.Middle, Spacing: layout.SpaceAround}
	// initialInset := layout.Inset{Top: unit.Dp(float32(newButton.sizeY))}

	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {

			area := PieChart(gtx.Ops, gtx, stats)
			return components.PieChartAreaButton.Layout(gtx, area)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.FinanceList(gtx.Ops, gtx, Ui)
		}),
	)
}
