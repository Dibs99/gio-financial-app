package components

import (
	"fmt"
	"gio-test/haslett/apiCalls"
	"gio-test/haslett/config"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func WeddingScreen(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	weddings := apiCalls.MyWeddings.Data.ReadHaslettWeddingss.Edges
	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H3(ui.theme, weddings[config.CurrentScreenID].Node.Name).Layout(gtx)
		}),
		layout.Rigid(weddingScreenChild(gtx, ui, "Date", weddings[config.CurrentScreenID].Node.Date)),
		layout.Rigid(weddingScreenChild(gtx, ui, "Stage", weddings[config.CurrentScreenID].Node.Stage)),
		layout.Rigid(weddingScreenChild(gtx, ui, "Package", fmt.Sprintf("%v", weddings[config.CurrentScreenID].Node.Package))),
	)
}

func weddingScreenChild(gtx layout.Context, ui *ThisUi, title string, value string) layout.Widget {
	myInset := unit.Dp(10)
	flex := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: myInset, Left: myInset, Right: myInset}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H6(ui.theme, title).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H6(ui.theme, value).Layout(gtx)
				}),
			)
		})
	}

}

func weddingText(gtx layout.Context, ui *ThisUi, wedding *apiCalls.WeddingNode) layout.Dimensions {
	flex := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
	bigText := material.Body1(ui.theme, fmt.Sprintf("%s\n%s", wedding.Node.Name, wedding.Node.Date))
	myInset := unit.Dp((100 - bigText.TextSize.V) / 2.7)
	return layout.Inset{Top: myInset, Left: myInset, Right: myInset}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return bigText.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(ui.theme, fmt.Sprintf("%v", wedding.Node.Package)).Layout(gtx)
			}),
		)
	})
}

func weddingBoxArea(gtx layout.Context, index *int) layout.Dimensions {
	const r = 10
	bounds := f32.Rect(0, 0, float32(gtx.Constraints.Min.X), float32(gtx.Constraints.Min.Y))
	defer op.Save(gtx.Ops).Load()
	clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Add(gtx.Ops)
	paint.ColorOp{Color: AreaButtonArray[*index].currentColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func weddingBox(gtx layout.Context, ui *ThisUi, wedding *apiCalls.WeddingNode, index *int) layout.Dimensions {
	gtx.Constraints.Min.Y = gtx.Px(unit.Dp(100))
	if config.CurrentScreenSize.X > config.MaxWidth {
		gtx.Constraints.Max.X = config.MaxWidth
		gtx.Constraints.Min.X = config.MaxWidth
	}

	area := weddingBoxArea(gtx, index)

	AreaButtonArray[*index].Layout(gtx, "wedding", *index, area)

	return weddingText(gtx, ui, wedding)
}

func WeddingList(gtx layout.Context, ui *ThisUi) layout.Dimensions {
	var array []layout.FlexChild

	for i, wedding := range apiCalls.MyWeddings.Data.ReadHaslettWeddingss.Edges {
		wedding := wedding
		i := i

		myFunc := func(gtx layout.Context) layout.Dimensions {
			return weddingBox(gtx, ui, &wedding, &i)
		}
		wrapper := layout.Rigid(myFunc)
		array = append(array, wrapper)
	}
	flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}

	return flex.Layout(gtx, array...)
	// return listWedding.Layout(gtx, len(array), func(gtx layout.Context, i int) layout.Dimensions {
	// 	return layout.UniformInset(unit.Dp(2)).Layout(gtx, array[i])
	// })
}

func weddingSetUpCallBack(gtx layout.Context) {
	// var startTime = time.Now()
	// var duration = 10 * time.Second

	// duration.Seconds()
	// time.Now().Sub(startTime)

	go func() {
		apiCalls.GetAllWeddings()
		for i := 0; i < len(apiCalls.MyWeddings.Data.ReadHaslettWeddingss.Edges); i++ {
			NewAreaButton()
		}
	}()

}
