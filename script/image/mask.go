package image

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
)

func (*Image) Mask(img, mask image.Image) (graph.Image, error) {
	return graph.Mask(img, mask)
}

// DrawMask draw img with mask over dest into a new image
func (*Image) DrawMask(img, mask, dest image.Image) (graph.Image, error) {
	return graph.DrawMask(img, mask, dest)
}
