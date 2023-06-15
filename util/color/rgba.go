package color

import "image/color"

func GreyScale(i int) color.Color {
	c := uint8(i & 0xff)
	return color.RGBA{R: c, G: c, B: c, A: 0xff}
}
