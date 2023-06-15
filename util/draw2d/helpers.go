package draw2d

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"image/color"
)

func FillPoly(gc *draw2dimg.GraphicContext, c color.Color, v ...float64) {
	gc.SetFillColor(c)
	gc.BeginPath()
	for i := 0; i < len(v); i += 2 {
		if i == 0 {
			gc.MoveTo(v[0], v[1])
		} else {
			gc.LineTo(v[i], v[i+1])
		}
	}
	gc.Close()
	gc.Fill()
}

func FillPolyRel(gc *draw2dimg.GraphicContext, c color.Color, v ...float64) {
	gc.SetFillColor(c)
	gc.BeginPath()
	var x, y float64
	for i := 0; i < len(v); i += 2 {
		x, y = x+v[i], y+v[i+1]
		if i == 0 {
			gc.MoveTo(x, y)
		} else {
			gc.LineTo(x, y)
		}
	}
	gc.Close()
	gc.Fill()
}

func FillRectangle(gc *draw2dimg.GraphicContext, x, y, w, h float64, c color.Color) (float64, float64) {
	gc.SetFillColor(c)
	gc.BeginPath()
	Rectangle(gc, x, y, w, h)
	gc.Fill()
	return x + w, y + h
}

func Rectangle(gc *draw2dimg.GraphicContext, x, y, w, h float64) {
	gc.MoveTo(x, y)
	gc.LineTo(x+w, y)
	gc.LineTo(x+w, y+h)
	gc.LineTo(x, y+h)
	gc.Close()
}

func RelLine(gc *draw2dimg.GraphicContext, x, y float64, v ...float64) {
	gc.MoveTo(x, y)
	for i := 0; i < len(v); i += 2 {
		x += v[i]
		y += v[i+1]
		gc.LineTo(x, y)
	}
}
