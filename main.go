package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	"gio-test/haslett/animate"
	"gio-test/haslett/apiCalls"
	"gio-test/haslett/components"
	"gio-test/haslett/config"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func main() {
	apiCalls.GetAllStats()
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(400), unit.Dp(500)),
			app.Title("Generic Financial App"),
		)

		if err := loop(w, apiCalls.MyStats); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()

}

var startTime = time.Now()
var duration = 10 * time.Second

func loop(w *app.Window, stats apiCalls.NewStats) error {
	Ui := components.NewUI()
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			config.Set(gtx, e.Size)
			var axis layout.Axis

			if float32(gtx.Constraints.Min.X/gtx.Constraints.Min.Y) < float32(0.6) {
				axis = layout.Vertical
			} else {
				axis = layout.Horizontal
			}
			masterWidgetList := []layout.Widget{}
			if apiCalls.Error != "" {
				masterWidgetList = append(masterWidgetList, func(gtx layout.Context) layout.Dimensions { return components.ErrorMessage(gtx, Ui) })
			}
			if config.CurrentScreen == "finance" {
				masterWidgetList = append(masterWidgetList, func(gtx layout.Context) layout.Dimensions { return financeLayout(axis, gtx, Ui) })

			}
			if config.CurrentScreen == "weddings" {
				masterWidgetList = append(masterWidgetList,
					func(gtx layout.Context) layout.Dimensions {
						if e.Size.X-config.MaxWidth > 0 {
							op.Offset(f32.Pt((float32(e.Size.X-config.MaxWidth) / 2), 0)).Add(gtx.Ops)
						}
						return components.WeddingList(gtx, Ui)
					})
			}
			if config.CurrentScreen == "wedding" {
				masterWidgetList = append(masterWidgetList,
					func(gtx layout.Context) layout.Dimensions {
						if e.Size.X-config.MaxWidth > 0 {
							op.Offset(f32.Pt((float32(e.Size.X-config.MaxWidth) / 2), 0)).Add(gtx.Ops)
						}
						return components.WeddingScreen(gtx, Ui)
					})
			}
			// layout the screen with space for footertabs at the bottom and space at the top for mobile
			masterList.Axis = layout.Vertical
			layout.Inset{Top: unit.Dp(40), Bottom: unit.Dp(float32(components.ButtonSizeY))}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return masterList.Layout(gtx, len(masterWidgetList), func(gtx layout.Context, i int) layout.Dimensions {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, masterWidgetList[i])
				})
			})
			animate.FadeInValue(gtx, config.Background, 254)
			// Push the footer tabs to the bottom of the screen
			footerSave := op.Save(gtx.Ops)
			op.Offset(f32.Pt(0, float32(gtx.Constraints.Max.Y-gtx.Px(unit.Dp(components.ButtonSizeY))))).Add(gtx.Ops)
			components.FooterTabs(gtx, e.Size, Ui)
			footerSave.Load()
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func drawImage(ops *op.Ops) {
	f, err := os.Open("./hamster.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err2 := png.Decode(f)
	if err2 != nil {
		log.Fatal(err2)
	}
	imageOp := paint.NewImageOp(img)
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4)))
	paint.PaintOp{}.Add(ops)
}

func PieChart(ops *op.Ops, gtx layout.Context) layout.Dimensions {
	var path clip.Path
	const r = 50 // roundness
	bounds := f32.Rect(0, 0, 100, 100)
	var variableCount float32 = 0
	var xCount float32 = 0
	var YCount float32 = 0

	for _, s := range apiCalls.MyStats.Data.ReadBankStatements {

		if s.MyCategory == "payments" {
			continue
		}
		value, err := strconv.ParseFloat(s.Percentage, 32)

		if err != nil {
			log.Fatal("PieChart, convert string to float32")
			log.Fatal(err)
		}
		percentage := float32(value)
		variable := (variableCount + percentage) * 4
		var X float32 = variable
		var Y float32 = 0
		// X1 and Y1 are the previous coordinates
		var X1 float32 = 0
		var Y1 float32 = 0
		// X2 and Y2 pull the line towards one of the corners to complete the path
		var X2 float32 = xCount
		var Y2 float32 = YCount
		// colour of this stat

		// X and Y coordinates are scaffold out via the variable
		if variable > 100 && variable < 200 {
			X = 100
			Y = variable - 100
			X1, Y1 = 100, 0
		}

		if variable >= 200 {
			X = 300 - variable
			Y = 100
			if variable >= 300 {
				X = 0
				Y = 400 - variable
				X1, Y1 = 0, 100
			} else {
				X1, Y1 = 100, 100
			}
		}

		if X2 > 100 && X2 < 200 {
			X2 = 100
			Y2 = X1 - 100
		}

		if X2 >= 200 {
			X2 = 100
			Y2 = 100
		}
		if X2 >= 300 {
			X2 = 0
			Y2 = 0
		}
		// if value < 0 {
		//  continue
		// }
		// used for debugging
		// fmt.Print(fmt.Sprintf("X: %v, Y: %v, X1: %v, Y1: %v, X2: %v, Y2: %v, Variable: %v\n", X, Y, X1, Y1, X2, Y2, variable))
		stack := op.Save(ops)
		path.Begin(ops)
		path.MoveTo(f32.Pt(50, 50))
		path.LineTo(f32.Pt(X, Y))
		path.LineTo(f32.Pt(X1, Y1))
		if percentage > 33 {
			path.LineTo(f32.Pt(Y1, Y1))
		}
		path.LineTo(f32.Pt(X2, Y2))
		path.LineTo(f32.Pt(50, 50))
		clip.Outline{Path: path.End()}.Op().Add(ops)
		clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Add(ops)

		paint.ColorOp{Color: s.Colour}.Add(ops)

		paint.PaintOp{}.Add(ops)
		stack.Load()
		variableCount += percentage
		xCount = X
		YCount = Y
	}
	return layout.Dimensions{Size: image.Point{100, 100}}
}

var (
	masterList = new(layout.List)
)

func newMasterList() *layout.List {
	newMasterList := new(layout.List)
	newMasterList.Axis = layout.Vertical

	return newMasterList
}
