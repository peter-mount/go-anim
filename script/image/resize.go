package image

import (
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/graph/resize"
	"image"
	"strings"
)

var (
	interpolationFunction = map[string]resize.InterpolationFunction{
		"nearestneighbor":   resize.NearestNeighbor,
		"bilinear":          resize.Bilinear,
		"bicubic":           resize.Bicubic,
		"mitchellnetravali": resize.MitchellNetravali,
		"lanczos2":          resize.Lanczos2,
		"lanczos3":          resize.Lanczos3,
	}
)

func getInterpolationFunction(interp string) resize.InterpolationFunction {
	interp = strings.ToLower(strings.TrimSpace(interp))
	if f, exists := interpolationFunction[interp]; exists {
		return f
	}
	return resize.NearestNeighbor
}

// Expand will add the specified number of pixels to the edges of an Image
func (_ *Image) Expand(src image.Image, top, left, bottom, right int) graph.Image {
	return graph.Expand(src, top, left, bottom, right)
}

func (_ *Image) Resize(width, height uint, img image.Image, interp string) image.Image {
	return resize.Resize(width, height, img, getInterpolationFunction(interp))
}

func (_ *Image) Thumbnail(maxWidth, maxHeight uint, img image.Image, interp string) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, img, getInterpolationFunction(interp))
}
