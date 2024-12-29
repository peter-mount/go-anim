package layout

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"image"
)

// Value is similar to Text except it consists of a Label as well as the text output.
// The label is rendered on the left of the component, whilst the value on the right
type Value struct {
	BaseComponent
	x, y           float64 // Position of text in component
	label          string  // Label for rendering on the left
	format         string  // Format for rendering on the right
	args           []any   // Args for format
	ll, lt, lr, lb float64 // copy of label string dimensions
	rl, rt, rr, rb float64 // copy of formatted string dimensions
}

func NewValue(label, format string, args ...any) *Value {
	t := &Value{
		label:         label,
		format:        format,
		args:          args,
		BaseComponent: BaseComponent{Type: "Value"},
	}
	t.BaseComponent.painter = t.paint
	return t
}

func (t *Value) Pos(x, y float64) *Value {
	t.x = x
	t.y = y
	t.updateRequired = true
	return t
}

func (t *Value) Args(args ...any) *Value {
	t.args = args
	t.updateRequired = true
	return t
}

func (t *Value) Layout(ctx draw2d.GraphicContext) bool {
	bounds := t.Bounds()

	t.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		lm, rm := t.metrics(gc)
		if bounds.Dx() == 0 {
			bounds.Max.X = bounds.Min.X + int(lm.MaxLineWidth) + int(rm.MaxLineWidth)
		}
		bounds.Max.Y = bounds.Min.Y + int(lm.MaxLineHeight*float64(max(len(lm.Lines), len(rm.Lines)))) + 2
		t.SetBounds(bounds)
	})

	t.updateRequired = false
	return true
}

func (t *Value) paint(gc *draw2dimg.GraphicContext) {
	lm, rm := t.metrics(gc)
	lm.Fill(gc)
	rm.Fill(gc)
}

func (t *Value) metrics(gc *draw2dimg.GraphicContext) (*util.AlignmentMetrics, *util.AlignmentMetrics) {
	bounds := t.Bounds()

	cx := bounds.Dx() >> 1

	offset := 4

	lm := util.RightAlignment.Metrics(gc, image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+cx-offset, bounds.Max.Y), 2, t.label)
	rm := util.LeftAlignment.Metrics(gc, image.Rect(bounds.Min.X+cx+offset, bounds.Min.Y, bounds.Max.X, bounds.Max.X), 2, t.format, t.args...)
	// Merge so they share some metric data ensuring they line up with each other
	lm.Merge(rm)

	return lm, rm
}

func (t *Value) String() string {
	return fmt.Sprintf(t.format, t.args...)
}
