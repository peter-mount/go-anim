package color

import (
	"image/color"
)

func Gradient(n int, from, to color.Color) []color.Color {
	fr, fg, fb, fa := from.RGBA()
	tr, tg, tb, ta := to.RGBA()
	dr, dg, db, da := diff(fr, tr, n), diff(fg, tg, n), diff(fb, tb, n), diff(fa, ta, n)

	var r []color.Color
	r = append(r, from)
	for i := 2; i < n; i++ {
		fr, fg, fb, fa = fr+dr, fg+dg, fb+db, fa+da
		r = append(r, color.RGBA{R: uint8(fr >> 8), G: uint8(fg >> 8), B: uint8(fb >> 8), A: uint8(fa >> 8)})
	}
	return append(r, to)
}

func diff(f, t uint32, n int) uint32 {
	if t < f {
		f, t = t, f
	}
	return (t - f) / uint32(n)
}
