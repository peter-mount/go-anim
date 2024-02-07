package exr

import (
	"image/color"
)

const (
	gammaFactor = 1.0 / 2.2
)

var (
	// RGBAModel returns the color.Model for RGBAColor colors.
	RGBAModel color.Model = color.ModelFunc(rgbaModel)
)

// RGBAColor represents a linear EXR color that implements the color.Color
// interface and is composed of R, G, B, and A components.
type RGBAColor struct {

	// R holds the amount of red in this color.
	R float32

	// G holds the amount of green in this color.
	G float32

	// B holds the amount of blue in this color.
	B float32

	// A holds the amount of alpha in this color.
	A float32
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the color. Each value ranges within [0, 0xffff], but is represented
// by a uint32 so that multiplying by a blend factor up to 0xffff will not
// overflow.
//
// An alpha-premultiplied color component c has been scaled by alpha (a),
// so has valid values 0 <= c <= a.
//
// Reinhard tone mapping and gamma correction are performed to convert the
// color into sRGB space.
func (c RGBAColor) RGBA() (uint32, uint32, uint32, uint32) {
	//// tone mapping
	//floatR := float64(c.R / (c.R + 1.0))
	//floatG := float64(c.G / (c.G + 1.0))
	//floatB := float64(c.B / (c.B + 1.0))
	//
	//// gamma correction
	//floatR = math.Pow(floatR, gammaFactor)
	//floatG = math.Pow(floatG, gammaFactor)
	//floatB = math.Pow(floatB, gammaFactor)
	//
	//// alpha pre-multiplication
	//floatR *= float64(c.A)
	//floatG *= float64(c.A)
	//floatB *= float64(c.A)

	const m = float32(0xFFFF)

	// uint32 conversion
	return uint32(c.R*m) & 0xFFFF,
		uint32(c.G*m) & 0xFFFF,
		uint32(c.B*m) & 0xFFFF,
		uint32(c.A*m) & 0xFFFF
}

func rgbaModel(c color.Color) color.Color {
	if _, ok := c.(RGBAColor); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	const m = float32(0xFFFF)
	return RGBAColor{
		R: float32(r) / m,
		G: float32(g) / m,
		B: float32(b) / m,
		A: float32(a) / m,
	}
}
