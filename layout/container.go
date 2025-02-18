package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
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
		if comp.IsUpdateRequired() && comp.Layout(ctx) {
			update = true
		}
	}
	return update
}

func (c *container) paint(gc *draw2dimg.GraphicContext) {
	for _, comp := range c.components {
		comp.Draw(gc)
	}
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
