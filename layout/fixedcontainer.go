package layout

import (
	"github.com/llgcode/draw2d"
	"image"
)

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
	//fmt.Println(c.GetType(), c.Bounds())
	if c.updateRequired {
		c.forceBounds()
	}

	if !c.container.Layout(ctx) {
		return false
	}

	c.container.Layout(ctx)
	return true
}
