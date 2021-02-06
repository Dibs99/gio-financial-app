package config

import (
	"image/color"
)

var (
	background    = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	Red           = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	Green         = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	Blue          = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
	CurrentScreen = "finance"
)
