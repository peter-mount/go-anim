package renderer

import (
	"errors"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/graph"
)

func CloneContext(ctx Context) Context {
	if c, ok := ctx.(*context); ok {
		n := &context{
			width:    c.width,
			height:   c.height,
			userdata: make(map[string]any),
		}

		// Copy the image as-is
		n.img = graph.DuplicateImage(c.Image())

		n.gc = draw2dimg.NewGraphicContext(n.img)

		for k, v := range c.userdata {
			n.userdata[k] = v
		}

		return n
	}

	panic(errors.New("not compatible context"))
}
