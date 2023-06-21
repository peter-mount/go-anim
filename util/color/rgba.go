package color

import "image/color"

func GreyScale(i int) color.Color {
	c := uint8(i & 0xff)
	return color.RGBA{R: c, G: c, B: c, A: 0xff}
}

// RGBA returns a color.RGBA instance based on the supplied 3 components
func RGBA(r, g, b uint32) color.RGBA {
	return color.RGBA{
		R: uint8((r >> 8) & 0xff),
		G: uint8((g >> 8) & 0xff),
		B: uint8((b >> 8) & 0xff),
		A: 255,
	}
}

// MinRGB returns the min value from the 3 components
func MinRGB(r, g, b uint32) uint32 {
	v := r
	if g < v {
		v = g
	}
	if b < v {
		v = b
	}
	return v
}

// MaxRGB returns the max value from the 3 components
func MaxRGB(r, g, b uint32) uint32 {
	v := r
	if g > v {
		v = g
	}
	if b > v {
		v = b
	}
	return v
}
