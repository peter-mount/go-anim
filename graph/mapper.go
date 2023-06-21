package graph

import (
	"image/color"
)

// Mapper is a function that can transform a color.Color to another
type Mapper func(col color.Color) (color.Color, error)

func (m Mapper) Filter() Filter {
	if m == nil {
		return nil
	}
	return func(_, _ int, col color.Color) (color.Color, error) {
		return m(col)
	}
}
