package config

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	background        = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: MyAlpha}
	Red               = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: MyAlpha}
	Green             = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: MyAlpha}
	Blue              = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: MyAlpha}
	White             = color.NRGBA{R: 228, G: 220, B: 220, A: MyAlpha}
	OffWhite          = color.NRGBA{R: 187, G: 169, B: 169, A: MyAlpha}
	MyAlpha           uint8
	CurrentScreen     = "finance"
	CurrentScreenID   int
	MaxWidth          int
	CurrentScreenSize image.Point
)

func Set(gtx layout.Context, e image.Point) {
	CurrentScreenSize = e
	MaxWidth = gtx.Px(unit.Dp(780))
}
