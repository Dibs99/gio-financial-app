package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var jsonString string = `{
    "data": [
        {
            "category" : "food",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "accommodation",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "travel",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "gifts",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        }
    ]
}`

type stat struct {
	Data []statChild `json:"data"`
}
type newStats struct {
	Data struct {
		ReadBankStatements []struct {
			Category   string      `json:"Category"`
			Percentage string      `json:"percentage"`
			Total      float32     `json:"total"`
			Colour     interface{} `json:"colour"`
		} `json:"readBankStatements"`
	} `json:"data"`
}
type statChild struct {
	Category   string      `json:"category"`
	Percentage float32     `json:"percentage"`
	Total      float32     `json:"total"`
	Colour     interface{} `json:"colour"`
}

func stats(jsonString string) *stat {
	// Declared an empty interface
	result := &stat{}
	myString := []byte(jsonString)
	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal(myString, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getStats() (record newStats) {
	postBody, _ := json.Marshal(map[string]string{
		"operationName": "HaslettBankStatements",
		"query":         "{readBankStatements{Category\nPercentage\nTotal}}",
	})
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8888/graphql", responseBody)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Origin", "http://fakewebsite.com")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("getStat: ")
		log.Fatal(err)
	}

	defer resp.Body.Close()

	err2 := json.NewDecoder(resp.Body).Decode(&record)
	// body, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		log.Fatal("getStat err2: ")
		log.Fatal(err2)
	}
	return record

}

func main() {
	// newStats :=
	// fmt.Print(d)
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(800), unit.Dp(400)),
			app.Title("Generic Financial App"),
		)

		if err := loop(w, &myStats); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, stats *newStats) error {
	ui := NewUI()
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			var axis layout.Axis
			if e.Size.X < 500 {
				axis = layout.Vertical
			} else {
				axis = layout.Horizontal
			}
			if currentScreen == "finance" {
				flex := layout.Flex{Axis: axis, Alignment: layout.Middle, Spacing: layout.SpaceAround}
				// initialInset := layout.Inset{Top: unit.Dp(float32(newButton.sizeY))}

				flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return PieChart(gtx.Ops, gtx, stats)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return financeList(gtx.Ops, gtx, ui)
					}),
				)
				footerTabs(gtx, e.Size.Y)
			} else {
				footerTabs(gtx, e.Size.Y)

			}

			e.Frame(gtx.Ops)

		}
	}
	return nil
}

type UI struct {
	theme  *material.Theme
	active int
	// list   layout.List
}

func NewUI() *UI {
	ui := &UI{
		theme: material.NewTheme(gofont.Collection()),
		// list: layout.List{
		//  Axis:      layout.Vertical,
		//  Alignment: layout.Middle,
		// },
	}

	return ui
}

type Button struct {
	pressed      bool
	currentColor color.NRGBA
	hoverColor   color.NRGBA
	initialColor color.NRGBA
	sizeX        int
	sizeY        int
}

func (b *Button) Layout(gtx layout.Context, screen string) layout.Dimensions {
	// Avoid affecting the input tree with pointer events.
	defer op.Push(gtx.Ops).Pop()

	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				currentScreen = screen
			case pointer.Enter:
				b.currentColor = b.hoverColor
			case pointer.Leave:
				b.currentColor = b.initialColor
			}
		}
	}

	// Confine the area for pointer events.
	pointer.Rect(image.Rect(0, 0, b.sizeX, b.sizeY)).Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)

	// Draw the button.
	return ColorBox(gtx, image.Pt(b.sizeX, b.sizeY), b.currentColor)
}

func PieChart(ops *op.Ops, gtx layout.Context, stats *newStats) layout.Dimensions {
	var path clip.Path
	const r = 50 // roundness
	bounds := f32.Rect(0, 0, 100, 100)
	var variableCount float32 = 0
	var xCount float32 = 0
	var YCount float32 = 0
	setColour := func(index int, colour interface{}) {
		stats.Data.ReadBankStatements[index].Colour = colour
	}
	for i, s := range stats.Data.ReadBankStatements {
		value, err := strconv.ParseFloat(s.Percentage, 32)
		if err != nil {
			log.Fatal("PiecChart, convert string to float32")
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

		random := func() uint8 {
			return uint8(rand.Intn(80) + 80)
		}

		if s.Colour == nil {
			colour := color.NRGBA{G: random(), B: random(), R: random(), A: 0xFF}
			s.Colour = colour
			setColour(i, colour)
		}
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

		// used for debugging
		// fmt.Print(fmt.Sprintf("X: %v, Y: %v, X1: %v, Y1: %v, X2: %v, Y2: %v, Variable: %v\n", X, Y, X1, Y1, X2, Y2, variable))
		stack := op.Push(ops)
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

		paint.ColorOp{Color: s.Colour.(color.NRGBA)}.Add(ops)

		paint.PaintOp{}.Add(ops)
		stack.Pop()
		variableCount += percentage
		xCount = X
		YCount = Y
	}
	return layout.Dimensions{Size: image.Point{100, 100}}
}

var (
	background    = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red           = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green         = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue          = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
	list          = &layout.List{Axis: layout.Vertical}
	masterList    = &layout.List{Axis: layout.Vertical}
	myStats       = getStats()
	button        = new(widget.Clickable)
	buttonState   = true
	financeButton = &Button{pressed: false, currentColor: blue, initialColor: blue, hoverColor: red, sizeX: 100, sizeY: 50}
	weddingButton = &Button{pressed: false, currentColor: green, initialColor: green, hoverColor: red, sizeX: 100, sizeY: 50}
	currentScreen = "finance"
)

func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
