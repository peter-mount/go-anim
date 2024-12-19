package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

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
	bounds := c.Bounds()

	c.FitToWidth()
	//fmt.Println("rct", c.Bounds())

	c.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		y := bounds.Min.Y
		for _, comp := range c.components {
			cb := comp.Bounds()
			dy := cb.Dy()
			cb.Min.Y = y
			cb.Max.Y = y + dy
			cb.Min.X = 0
			cb.Max.X = bounds.Dx()
			comp.SetBounds(cb)
			//fmt.Println("rct", comp.GetType(), comp.Bounds())
			y = y + dy
		}
		bounds.Max.Y = y
		c.SetBounds(bounds)

		c.container.Layout(gc)
	})

	c.FitToHeight()

	return true
}
