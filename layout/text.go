package layout

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"math"
)

type Text struct {
	BaseComponent
	x, y         float64 // Position of text in component
	format       string
	args         []any
	l, top, r, b float64 // copy of string dimensions
}

func NewText(format string, args ...any) *Text {
	t := &Text{
		format:        format,
		args:          args,
		BaseComponent: BaseComponent{Type: "Text"},
	}
	t.painter = t.paint
	return t
}

func (t *Text) Pos(x, y float64) *Text {
	t.x = x
	t.y = y
	t.updateRequired = true
	return t
}

func (t *Text) Args(args ...any) *Text {
	t.args = args
	t.updateRequired = true
	return t
}

func (t *Text) Layout(ctx draw2d.GraphicContext) bool {
	bounds := t.Bounds()

	t.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		t.l, t.top, t.r, t.b = ctx.GetStringBounds(t.String())
		if bounds.Dx() == 0 {
			bounds.Max.X = bounds.Min.X + int(t.r-t.l)
		}
		bounds.Max.Y = bounds.Min.Y + int(-math.Floor(t.top))
		t.SetBounds(bounds)
	})

	t.updateRequired = false
	return true
}

func (t *Text) paint(gc *draw2dimg.GraphicContext) {
	bounds := t.Bounds()
	s := t.String()
	switch t.alignment {
	case LeftAlignment:
		util.DrawStringLeft(gc, 0, float64(bounds.Dy())-t.b, s)
	case CenterAlignment:
		util.DrawStringCenter(gc, float64(bounds.Dx())-t.b, float64(bounds.Dy())/2, s)
	case RightAlignment:
		util.DrawStringRight(gc, float64(bounds.Dx()), float64(bounds.Dy())-t.b, s)
	}
}

func (t *Text) String() string {
	return fmt.Sprintf(t.format, t.args...)
}
