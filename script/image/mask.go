package image

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
	"image/color"
)

func (g *Image) Mask(img, mask image.Image, r0, g0, b0 uint32) (graph.Image, error) {
	return g.FilterNew(
		func(x, y int, col color.Color) (color.Color, error) {
			r1, g1, b1, a1 := mask.At(x, y).RGBA()
			if a1 == 0 || (r1 != r0 && g1 != g0 && b1 != b0) {
				return col, nil
			}
			return color.Black, nil
		},
		img)
}
