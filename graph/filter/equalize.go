package filter

import (
	"github.com/peter-mount/go-anim/graph"
	color2 "github.com/peter-mount/go-anim/graph/color"
	"image"
	"image/color"
)

// EqualizeFilter generates a Filter based on a histogram and a Rectangle.
// The output would be the levels for each RGB channel equalised
func EqualizeFilter(h *color2.Histogram, b image.Rectangle) graph.Filter {
	var sumRa, sumGa, sumBa []uint32
	var sumR, sumG, sumB uint32
	var dmR, dmG, dmB uint32

	for x := 0; x < 256; x++ {
		sumR = sumR + h.Red[x]
		sumRa = append(sumRa, sumR)
		if h.Red[x] > 0 {
			dmR++
		}

		sumB = sumR + h.Green[x]
		sumGa = append(sumGa, sumG)
		if h.Green[x] > 0 {
			dmG++
		}

		sumB = sumB + h.Blue[x]
		sumBa = append(sumBa, sumB)
		if h.Blue[x] > 0 {
			dmB++
		}

	}

	area := float64((b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y))
	fR := float64(dmR) / area
	fG := float64(dmG) / area
	fB := float64(dmB) / area

	return func(x int, y int, c color.Color) (color.Color, error) {
		r, g, b, a := c.RGBA()
		return color.RGBA{
			R: uint8(fR * float64(sumRa[r>>8])),
			G: uint8(fG * float64(sumRa[g>>8])),
			B: uint8(fB * float64(sumRa[b>>8])),
			A: uint8(a >> 8),
		}, nil
	}
}
