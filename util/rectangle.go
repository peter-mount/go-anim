package util

import (
	"github.com/llgcode/draw2d"
	"image"
	"math"
	"strings"
)

type Rectangle struct {
	X1, Y1, X2, Y2 float64
}

func Rect(x1, y1, x2, y2 float64) Rectangle {
	return Rectangle{X1: math.Min(x1, x2), Y1: math.Min(y1, y2), X2: math.Max(x1, x2), Y2: math.Max(y1, y2)}
}

func (r Rectangle) Width() float64 { return r.X2 - r.X1 }

func (r Rectangle) Height() float64 { return r.Y2 - r.Y1 }

// Add returns a Rectangle that contains both rectangles
func (r Rectangle) Add(b Rectangle) Rectangle {
	return Rect(math.Min(r.X1, b.X1), math.Min(r.Y1, b.Y1), math.Max(r.X2, b.X2), math.Max(r.Y2, b.Y2))
}

func (r Rectangle) IsZero() bool {
	return r.X1 == 0 && r.Y1 == 0 && r.X2 == 0 && r.Y2 == 0
}

func (r Rectangle) Expand(l, t, right, b float64) Rectangle {
	return Rect(r.X1-l, r.Y1-t, r.X2+right, r.Y2+b)
}

func (r Rectangle) Reduce(l, t, right, b float64) Rectangle {
	return Rect(r.X1+l, r.Y1+t, r.X2-right, r.Y2-b)
}

func (r Rectangle) String() string {
	if r.IsZero() {
		return ""
	}
	return strings.Join([]string{
		FloatToA(r.X1),
		FloatToA(r.Y1),
		FloatToA(r.X2),
		FloatToA(r.Y2),
	}, ",")
}

func (r Rectangle) Rect() image.Rectangle {
	return image.Rect(int(r.X1), int(r.Y1), int(r.X2), int(r.Y2))
}

func (r Rectangle) AddPath(ctx draw2d.GraphicContext) {
	ctx.MoveTo(r.X1, r.Y1)
	ctx.LineTo(r.X2, r.Y1)
	ctx.LineTo(r.X2, r.Y2)
	ctx.LineTo(r.X1, r.Y2)
	// Don't close as start might not be this rectangle
	ctx.LineTo(r.X1, r.Y1)
}

func RectFromRect(rect image.Rectangle) Rectangle {
	return Rect(float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y))
}
