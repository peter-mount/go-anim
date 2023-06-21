package color

import (
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/util"
	"image/color"
)

// MinLevel is a Mapper that will ensure the image contains the minimum of each component.
// This can be used to set the upper limit in an image, clipping entries that are too bright.
func MinLevel(r1, g1, b1 uint32) graph.Mapper {
	return Level(r1, g1, b1, util.MinU32)
}

// MaxLevel is a Mapper that will ensure the image contains the maximum of each component.
// This can be used to set the lower limit in an image, clipping entries that are too dark.
func MaxLevel(r1, g1, b1 uint32) graph.Mapper {
	return Level(r1, g1, b1, util.MaxU32)
}

// Level will apply a function against the R, G and B components of a pixel, replacing it with the result.
// This is the implementation behind MinLevel and MaxLevel which use the util.MinU32 and util.MaxU32
// functions.
// The first parameter passed to this function are the parameters of this function.
// The second parameter is from the pixel.
func Level(r1, g1, b1 uint32, f func(a, b uint32) uint32) graph.Mapper {
	return func(col color.Color) (color.Color, error) {
		r2, g2, b2, a := col.RGBA()
		return color.RGBA{
			R: uint8(f(r1, r2) >> 8),
			G: uint8(f(g1, g2) >> 8),
			B: uint8(f(b1, b2) >> 8),
			A: uint8(a),
		}, nil
	}
}

// Brighten is a Mapper that will increase the R, G and B components by a fixed amount.
func Brighten(amount uint32) graph.Mapper {
	return func(col color.Color) (color.Color, error) {
		r, g, b, a := col.RGBA()
		return color.RGBA{
			R: uint8((r + amount) >> 8),
			G: uint8((g + amount) >> 8),
			B: uint8((b + amount) >> 8),
			A: uint8(a),
		}, nil
	}
}

// Darken is a Mapper that will decrease the R, G and B components by a fixed amount.
func Darken(amount uint32) graph.Mapper {
	return func(col color.Color) (color.Color, error) {
		r, g, b, a := col.RGBA()
		return color.RGBA{
			R: uint8((r - amount) >> 8),
			G: uint8((g - amount) >> 8),
			B: uint8((b - amount) >> 8),
			A: uint8(a),
		}, nil
	}
}
