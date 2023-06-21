package graph

import (
	"image"
	"image/color"
)

// Filter is a basic mapping function that is supplied with a Pixel and returns the new colour for that Pixel.
type Filter func(x, y int, col color.Color) (color.Color, error)

// Of returns a Filter based on the supplied filters
func Of(filters ...Filter) Filter {
	var f Filter
	for _, filter := range filters {
		f = f.Then(filter)
	}
	return f
}

// Then returns a Filter that will run this one then the subsequent one in sequence
func (f Filter) Then(b Filter) Filter {
	if f == nil {
		return b
	}
	if b == nil {
		return f
	}
	return func(x int, y int, c color.Color) (color.Color, error) {
		nc, err := f(x, y, c)
		if err != nil {
			return nil, err
		}
		return b(x, y, nc)
	}
}

// Do apply the filter against the source image, writing the result to the destination
// over the area defined by the supplied rectangle.
func (f Filter) Do(src image.Image, dst Image, b image.Rectangle) error {
	if f == nil {
		return nil
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c, err := f(x, y, src.At(x, y))

			if err != nil {
				return err
			}

			if dst != nil {
				dst.Set(x, y, c)
			}
		}
	}

	return nil
}

// DoNew applies the filter against an image, returning a new mutable image with the result.
func (f Filter) DoNew(src image.Image) (Image, error) {
	dst := NewRGBA(src)
	err := f.Do(src, dst, dst.Bounds())
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// DoOver applies the filter over the supplied mutable image, overwriting its previous state.
func (f Filter) DoOver(src Image) error {
	return f.Do(src, src, src.Bounds())
}
