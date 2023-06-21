package color

import (
	"github.com/peter-mount/go-anim/graph"
	"image/color"
)

// Map applies a simple mapping function against the pixel Map.
func Map(f graph.Mapper) graph.Filter {
	return func(_, _ int, col color.Color) (color.Color, error) {
		return f(col)
	}
}

// InvertColor is a Mapper to invert an image
func InvertColor(col color.Color) color.Color {
	r, g, b, a := col.RGBA()
	return color.RGBA{
		R: uint8((0xFFFFFFFF - r) >> 8),
		G: uint8((0xFFFFFFFF - g) >> 8),
		B: uint8((0xFFFFFFFF - b) >> 8),
		A: uint8(a),
	}
}
