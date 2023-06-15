package script

import (
	color2 "github.com/peter-mount/go-anim/util/color"
	"image/color"
)

type Colour struct{}

func (_ Colour) GreyScale(i int) color.Color {
	return color2.GreyScale(i)
}

func (_ Colour) ColourString(c color.Color) string {
	return color2.ColourString(c)
}

func (_ Colour) Colour(hex string) (color.RGBA, error) {
	return color2.ParseColour(hex)
}

func (_ Colour) Gradient(n int, from, to color.Color) []color.Color {
	return color2.Gradient(n, from, to)
}
