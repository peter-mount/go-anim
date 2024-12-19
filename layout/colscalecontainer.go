package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

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

	//fmt.Printf("\nLayout %q\t%v\n", c.Type, c.scales)

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

				cb.Min.X = bounds.Min.X
				cb.Min.Y = 0
				cb.Max.X = cb.Min.X + int(width*scale)
				cb.Max.Y = bounds.Dy()

				comp.SetBounds(cb)

				// Move to next component
				bounds.Min.X = cb.Max.X
			}
		}

		c.container.Layout(gc)
	})

	c.FitToHeight()

	return true
}
