package image

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
)

func (g *Image) Mask(img, mask image.Image) (graph.Image, error) {
	return graph.Mask(img, mask)
}
