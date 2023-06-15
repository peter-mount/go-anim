package unit

import (
	"fmt"
	"git.area51.dev/peter/videoident/util"
	"github.com/llgcode/draw2d/draw2dimg"
	"strings"
)

type Dimension struct {
	Top    Value
	Bottom Value
	Left   Value
	Right  Value
}

func SquareDimension(v Value) Dimension {
	return NewDimension(v, v, v, v)
}

func RectangleDimension(t, r Value) Dimension {
	return NewDimension(t, t, r, r)
}

func NewDimension(t, b, l, r Value) Dimension {
	return Dimension{Top: t, Bottom: b, Left: l, Right: r}
}

func (d Dimension) IsZero() bool {
	return d.Top.IsZero() && d.Bottom.IsZero() && d.Left.IsZero() && d.Right.IsZero()
}

func (d Dimension) Convert(gc *draw2dimg.GraphicContext, to Unit) Dimension {
	return Dimension{
		Top:    d.Top.Convert(gc, to),
		Bottom: d.Bottom.Convert(gc, to),
		Left:   d.Left.Convert(gc, to),
		Right:  d.Right.Convert(gc, to),
	}
}

func (d Dimension) String() string {
	if d.IsZero() {
		return ""
	}
	tbe, lre := d.Top.Equals(d.Bottom), d.Left.Equals(d.Right)
	switch {
	case tbe && lre && d.Top.Equals(d.Right):
		return d.Top.String()

	case tbe && lre:
		return d.Top.String() + " " + d.Right.String()

	case !tbe && lre:
		return d.Top.String() + " " + d.Right.String() + " " + d.Bottom.String()

	default:
		return d.Top.String() + " " + d.Right.String() + " " + d.Bottom.String() + " " + d.Left.String()
	}
}

// ParseDimension parses a string consisting of 1 to 4 Values that describe a Dimension.
// When 1 value is specified, it applies to all 4 sides
// When 2 values are specified, first is for Top and Bottom, second for Left & Right
// When 3 values, first is top, second left & right, and third bottom
// When 4 then in this order: top,right,bottom,left
//
// This is based on how css works going clockwise from the top
// https://developer.mozilla.org/en-US/docs/Web/CSS/margin
func ParseDimension(s string) (Dimension, error) {
	var a []Value
	for _, s1 := range strings.Split(s, " ") {
		s1 = strings.TrimSpace(s1)
		if s1 != "" {
			v, err := ParseValue(s1)
			if err != nil {
				return Dimension{}, err
			}
			a = append(a, v)
		}
	}

	switch len(a) {
	case 0:
		return Dimension{}, nil

	case 1:
		return Dimension{
			Top:    a[0],
			Bottom: a[0],
			Left:   a[0],
			Right:  a[0],
		}, nil

	case 2:
		return Dimension{
			Top:    a[0],
			Bottom: a[0],
			Left:   a[1],
			Right:  a[1],
		}, nil

	case 3:
		return Dimension{
			Top:    a[0],
			Left:   a[1],
			Right:  a[1],
			Bottom: a[2],
		}, nil

	case 4:
		return Dimension{
			Top:    a[1],
			Right:  a[2],
			Bottom: a[3],
			Left:   a[4],
		}, nil

	default:
		return Dimension{}, fmt.Errorf("rectangle %q must be 1,2,3 or 4 values", s)
	}
}

// Add adds two dimensions, increasing the values for each side.
// The result will be of the dimension on the left hand side
func (d Dimension) Add(gc *draw2dimg.GraphicContext, b Dimension) Dimension {
	return Dimension{
		Top:    d.Top.Add(gc, b.Top),
		Bottom: d.Bottom.Add(gc, b.Bottom),
		Left:   d.Left.Add(gc, b.Left),
		Right:  d.Right.Add(gc, b.Right),
	}
}

// Sub subtracts two dimensions, increasing the values for each side.
// The result will be of the dimension on the left hand side
func (d Dimension) Sub(gc *draw2dimg.GraphicContext, b Dimension) Dimension {
	return Dimension{
		Top:    d.Top.Sub(gc, b.Top),
		Bottom: d.Bottom.Sub(gc, b.Bottom),
		Left:   d.Left.Sub(gc, b.Left),
		Right:  d.Right.Sub(gc, b.Right),
	}
}

// ReduceRect reduces a rectangle by the dimension
func (d Dimension) ReduceRect(gc *draw2dimg.GraphicContext, r util.Rectangle) util.Rectangle {
	return r.Reduce(d.Left.Pixels(gc), d.Top.Pixels(gc), d.Right.Pixels(gc), d.Bottom.Pixels(gc))
}

// ExpandRect expands a rectangle by the dimension
func (d Dimension) ExpandRect(gc *draw2dimg.GraphicContext, r util.Rectangle) util.Rectangle {
	return r.Expand(d.Left.Pixels(gc), d.Top.Pixels(gc), d.Right.Pixels(gc), d.Bottom.Pixels(gc))
}
