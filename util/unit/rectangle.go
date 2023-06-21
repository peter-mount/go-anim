package unit

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"strings"
)

type Rectangle struct {
	X1 Value
	Y1 Value
	X2 Value
	Y2 Value
}

func (d Rectangle) IsZero() bool {
	return d.X1.IsZero() && d.Y1.IsZero() && d.X2.IsZero() && d.Y2.IsZero()
}
func (d Rectangle) Convert(gc *draw2dimg.GraphicContext, to Unit) Rectangle {
	return Rectangle{
		X1: d.X1.Convert(gc, to),
		Y1: d.Y1.Convert(gc, to),
		X2: d.X2.Convert(gc, to),
		Y2: d.Y2.Convert(gc, to),
	}
}

func (d Rectangle) Height(gc *draw2dimg.GraphicContext) Value {
	return d.Y2.Sub(gc, d.Y1)
}

func (d Rectangle) Width(gc *draw2dimg.GraphicContext) Value {
	return d.X2.Sub(gc, d.X1)
}

func (d Rectangle) Add(gc *draw2dimg.GraphicContext, b Rectangle) Rectangle {
	return Rectangle{
		X1: d.X1.Min(gc, b.X1),
		Y1: d.Y1.Min(gc, b.Y1),
		X2: d.X2.Max(gc, b.X2),
		Y2: d.Y2.Max(gc, b.Y2),
	}
}

func (d Rectangle) Expand(gc *draw2dimg.GraphicContext, dim Dimension) Rectangle {
	return Rectangle{
		X1: d.X1.Sub(gc, dim.Left),
		Y1: d.Y1.Sub(gc, dim.Top),
		X2: d.X2.Add(gc, dim.Right),
		Y2: d.Y2.Add(gc, dim.Bottom),
	}
}

func (d Rectangle) Reduce(gc *draw2dimg.GraphicContext, dim Dimension) Rectangle {
	return Rectangle{
		X1: d.X1.Add(gc, dim.Left),
		Y1: d.Y1.Add(gc, dim.Top),
		X2: d.X2.Sub(gc, dim.Right),
		Y2: d.Y2.Sub(gc, dim.Bottom),
	}
}

func (d Rectangle) Rectangle(gc *draw2dimg.GraphicContext) util.Rectangle {
	return util.Rect(d.X1.Pixels(gc), d.Y1.Pixels(gc), d.X2.Pixels(gc), d.Y2.Pixels(gc))
}

// FromRectangle takes a plain rectangle and converts it into one with the specified units
func FromRectangle(gc *draw2dimg.GraphicContext, r util.Rectangle, unit Unit) Rectangle {
	d := Rectangle{X1: Pixels(r.X1), Y1: Pixels(r.Y1), X2: Pixels(r.X2), Y2: Pixels(r.Y2)}
	return d.Convert(gc, unit)
}

func (d Rectangle) String() string {
	if d.IsZero() {
		return ""
	}
	return d.X1.String() + " " + d.Y1.String() + " " + d.X2.String() + " " + d.Y2.String()
}

func (d Rectangle) SizeString() string {
	if d.IsZero() {
		return ""
	}
	return d.X1.String() + " " + d.Y1.String() + " " + d.X2.Sub(nil, d.X1).String() + " " + d.Y2.Sub(nil, d.Y1).String()
}

// ParseRectangle parses a string consisting of 4 Values that describe a Rectangle
func ParseRectangle(s string) (Rectangle, error) {
	var a []Value
	var as []string
	for _, s1 := range strings.Split(s, " ") {
		s1 = strings.TrimSpace(s1)
		if s1 != "" {
			v, err := ParseValue(s1)
			if err != nil {
				return Rectangle{}, err
			}
			a = append(a, v)
			as = append(as, s1)
		}
	}

	switch len(a) {
	case 0:
		return Rectangle{}, nil

	case 4:
		// FIXME this has no verification on the ordering of the rectangle
		// e.g. right now X1>X2 or Y1>Y2 are possible
		return Rectangle{X1: a[0], Y1: a[1], X2: a[2], Y2: a[3]}, nil

	default:
		return Rectangle{}, fmt.Errorf("rectangle %q must be 4 values", s)
	}
}

func ParseSizedRectangle(s string) (Rectangle, error) {
	r, err := ParseRectangle(s)
	if err == nil {
		r.X2 = r.X1.Add(nil, r.X2)
		r.Y2 = r.Y1.Add(nil, r.Y2)
	}
	return r, err
}
