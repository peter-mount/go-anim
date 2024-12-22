package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
)

// RowContainer contains components vertically.
func RowContainer() Container {
	return newRowContainer()
}

func newRowContainer() *rowContainer {
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
	bounds := c.Bounds()

	c.FitToWidth()

	c.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		y := c.insetMinY
		for _, comp := range c.components {
			comp.Layout(gc)
			cb := comp.Bounds()
			dy := cb.Dy()
			comp.SetBounds(image.Rect(c.insetX, y, bounds.Dx()-c.insetX, y+dy))
			comp.Layout(gc)
			y = y + comp.Bounds().Dy()
		}
		bounds.Max.Y = bounds.Min.Y + y + c.insetMaxY
		c.SetBounds(bounds)

		c.container.Layout(gc)
	})

	return true
}
