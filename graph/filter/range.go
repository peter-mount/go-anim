package filter

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
	"image/color"
)

// ClipInside will change all pixels outside the Rectangle to be transparent
func ClipInside(rect image.Rectangle) graph.Filter {
	return clipRect(rect, true)
}

// ClipOutside will change all pixels inside the Rectangle to be transparent
func ClipOutside(rect image.Rectangle) graph.Filter {
	return clipRect(rect, false)
}

func clipRect(rect image.Rectangle, inside bool) graph.Filter {
	return func(x, y int, col color.Color) (color.Color, error) {
		p := image.Point{X: x, Y: y}
		if p.In(rect) == inside {
			return col, nil
		}
		return color.Transparent, nil
	}
}
