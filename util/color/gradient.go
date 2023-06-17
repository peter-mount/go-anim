package color

import (
	"image/color"
)

func Gradient(n int, from, to color.Color) []color.Color {
	fr, fg, fb, fa := from.RGBA()
	tr, tg, tb, ta := to.RGBA()
	dr, dg, db, da := diff(fr, tr, n), diff(fg, tg, n), diff(fb, tb, n), diff(fa, ta, n)
	cr, cg, cb, ca := int32(fr), int32(fg), int32(fb), int32(fa)

	var r []color.Color
	r = append(r, from)
	for i := 2; i < n; i++ {
		cr, cg, cb, ca = cr+dr, cg+dg, cb+db, ca+da
		r = append(r, color.RGBA{R: uint8(cr >> 8), G: uint8(cg >> 8), B: uint8(cb >> 8), A: uint8(ca >> 8)})
	}

	r = append(r, to)

	return r
}

func diff(f, t uint32, n int) int32 {
	/*if t < f {
		f, t = t, f
	}*/
	return (int32(t) - int32(f)) / int32(n-1)
}
