package image

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
)

// Crop will crop an image to fit the required bounds
func (_ *Image) Crop(src image.Image, bounds image.Rectangle) graph.Image {
	return graph.Crop(src, bounds)
}

func (_ *Image) AutoCrop(src image.Image) graph.Image {
	return graph.AutoCrop(src)
}
