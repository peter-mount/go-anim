package filter

import (
	"image/color"
)

// BlackMaskFilter returns a filter which will mark a pixel as Opaque if it's black
// otherwise leaves the colour as is
func BlackMaskFilter(_, _ int, col color.Color) (color.Color, error) {
	return maskFilter(0, 0, 0, col)
}

// WhiteMaskFilter returns a filter which will mark a pixel as Opaque if it's white
// otherwise Transparent
func WhiteMaskFilter(_, _ int, col color.Color) (color.Color, error) {
	return maskFilter(0xffff, 0xffff, 0xffff, col)
}

func maskFilter(r1, g1, b1 uint32, col color.Color) (color.Color, error) {
	r, g, b, _ := col.RGBA()
	if r == r1 && g == g1 && b == b1 {
		return color.Opaque, nil
	}
	return col, nil
}
