package layout

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Text is a simple Component which renders a formatted string.
type Text struct {
	BaseComponent
	x, y         float64 // Position of text in component
	format       string
	args         []any
	l, top, r, b float64 // copy of string dimensions
}

// NewText creates a new Text component using the provided string. The string is passed to fmt.Sprintf() so it uses the
// same formatting.
func NewText(format string, args ...any) *Text {
	t := &Text{
		format:        format,
		args:          args,
		BaseComponent: BaseComponent{Type: "Text"},
	}
	t.BaseComponent.painter = t.paint
	return t
}

func (t *Text) Pos(x, y float64) *Text {
	t.x = x
	t.y = y
	t.updateRequired = true
	return t
}

// Args allows the args passed to NewText to be replaced, causing the component's output to change
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
		bounds.Max.Y = bounds.Min.Y + int(t.b-t.top) + 4
		t.SetBounds(bounds)
	})

	t.updateRequired = false
	return true
}

func (t *Text) paint(gc *draw2dimg.GraphicContext) {
	t.alignment.Fill(gc, t.Bounds(), 2, t.format, t.args...)
}

func (t *Text) String() string {
	return fmt.Sprintf(t.format, t.args...)
}
