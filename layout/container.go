package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"golang.org/x/image/colornames"
	"image"
)

type Container interface {
	Component
	Add(Component) Container
	IsEmpty() bool
}

type container struct {
	BaseComponent
	components []Component // Components in the container
}

func (c *container) Width() int {
	w := 0
	for _, component := range c.components {
		w += component.Width()
	}
	return w
}

func (c *container) Height() int {
	h := 0
	for _, component := range c.components {
		h += component.Height()
	}
	return h
}

func (c *container) Add(comp Component) Container {
	c.components = append(c.components, comp)
	c.updateRequired = true
	return c
}

func (c *container) IsEmpty() bool {
	return len(c.components) == 0
}

func (c *container) Layout(ctx draw2d.GraphicContext) bool {
	update := false
	for _, comp := range c.components {
		if comp.Layout(ctx) {
			update = true
		}
	}
	return update
}

func (c *container) paint(gc *draw2dimg.GraphicContext) {
	for _, comp := range c.components {
		comp.Draw(gc)
	}
	gc.SetStrokeColor(colornames.Red)
	gc.BeginPath()
	util.RectFromRect(c.Bounds()).AddPath(gc)
	gc.Stroke()
}

func (c *container) FitToHeight() int {
	height := 0
	for _, comp := range c.components {
		cb := comp.Bounds()
		if cb.Dy() > height {
			height = cb.Dy()
		}
	}
	b := c.Bounds()
	b.Max.Y = b.Min.Y + height
	c.SetBounds(b)
	return height
}

func (c *container) FitToWidth() int {
	width := 0
	for _, comp := range c.components {
		cb := comp.Bounds()
		if cb.Dx() > width {
			width = cb.Dx()
		}
	}
	b := c.Bounds()
	b.Max.X = b.Min.X + width
	c.SetBounds(b)
	return width
}

func FixedContainer(rect image.Rectangle) Container {
	c := &fixedContainer{
		container: container{
			BaseComponent: BaseComponent{Type: "FixedContainer"},
		},
	}
	c.SetBounds(rect)
	c.BaseComponent.painter = c.paint
	c.forceBounds()
	return c
}

type fixedContainer struct {
	container
}

func (c *fixedContainer) forceBounds() {
	for _, comp := range c.components {
		comp.SetBounds(c.bounds)
	}
}

func (c *fixedContainer) Layout(ctx draw2d.GraphicContext) bool {
	if c.updateRequired {
		c.forceBounds()
	}

	if !c.container.Layout(ctx) {
		return false
	}

	c.container.Layout(ctx)
	return true
}

// ColScaleContainer lays out it's components based on specific scaling horizontally
func ColScaleContainer(scales ...float64) Container {
	c := &colScaleContainer{
		container: container{
			BaseComponent: BaseComponent{Type: "ColScaleContainer"},
		},
		scales: scales,
	}
	c.BaseComponent.painter = c.paint
	return c
}

type colScaleContainer struct {
	container
	scales []float64
}

func (c *colScaleContainer) Layout(ctx draw2d.GraphicContext) bool {
	if !c.container.Layout(ctx) {
		return false
	}

	// Get max height of this row
	c.FitToHeight()

	bounds := c.Bounds()
	c.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		// Now update the widths of this row
		width := float64(bounds.Dx())
		for i, scale := range c.scales {
			if i < len(c.components) {
				comp := c.components[i]

				comp.Layout(gc)

				cb := comp.Bounds()
				cb.Min = bounds.Min
				cb.Max.X = cb.Min.X + int(width*scale)
				cb.Max.Y = bounds.Max.Y

				comp.SetBounds(cb)
				comp.Layout(gc)

				// Move to next component
				bounds.Min.X = cb.Max.X
			}
		}

		c.container.Layout(gc)
	})

	c.FitToHeight()

	return true
}

// RowContainer contains components vertically.
func RowContainer() Container {
	c := &rowContainer{
		container: container{
			BaseComponent: BaseComponent{Type: "RowContainer"},
		},
	}
	c.BaseComponent.painter = c.paint
	return c
}

type rowContainer struct {
	container
}

func (c *rowContainer) Layout(ctx draw2d.GraphicContext) bool {
	if !c.container.Layout(ctx) && !c.updateRequired {
		return false
	}

	bounds := c.Bounds()

	c.FitToWidth()

	c.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		y := bounds.Min.Y
		for _, comp := range c.components {
			cb := comp.Bounds()
			dy := cb.Dy()
			cb.Min.Y = y
			cb.Max.Y = y + dy
			cb.Min.X = bounds.Min.X
			cb.Max.X = bounds.Max.X
			comp.SetBounds(cb)
			comp.Layout(gc)
			y = y + dy
		}
		bounds.Max.Y = y
		c.SetBounds(bounds)

		c.container.Layout(gc)
	})

	c.FitToHeight()

	return true
}
