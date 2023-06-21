package util

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util/draw2d"
	"image/color"
)

func DrawColourBars(gc *draw2dimg.GraphicContext, bounds Rectangle, cols ...color.Color) (float64, float64) {
	x, y, w, h := bounds.X1, bounds.Y1, bounds.Width(), bounds.Height()
	l := float64(len(cols))
	dw := w / l
	for _, col := range cols {
		x, _ = draw2d.FillRectangle(gc, x, y, dw, h, col)
	}
	return dw, l
}

func DrawColourBarsVertical(gc *draw2dimg.GraphicContext, bounds Rectangle, cols ...color.Color) (float64, float64) {
	x, y, w, h := bounds.X1, bounds.Y1, bounds.Width(), bounds.Height()
	l := float64(len(cols))
	dh := h / l
	for _, col := range cols {
		_, y = draw2d.FillRectangle(gc, x, y, w, dh, col)
	}
	return dh, l
}
