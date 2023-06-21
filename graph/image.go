// Package graph is a suite of utilities for performing various operation's on an Image
package graph

import (
	"github.com/peter-mount/go-anim/util"
	"image"
	"image/color"
	"image/draw"
)

// Image is an image.Image with a Set method to change a single pixel.
type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

// NewRGBA creates a new mutable image with the same dimensions of another image
func NewRGBA(img image.Image) Image {
	return NewRGBAImage(img.Bounds())
}

func NewRGBAImage(bounds image.Rectangle) Image {
	return image.NewRGBA(bounds)
}

// DuplicateImage creates a new copy of an image which is also mutable
func DuplicateImage(img image.Image) Image {
	dst := NewRGBA(img)
	draw.Draw(dst, img.Bounds(), img, image.Point{}, draw.Src)
	return dst
}

type wrapper struct {
	img image.Image
}

func (w *wrapper) ColorModel() color.Model     { return w.img.ColorModel() }
func (w *wrapper) Bounds() image.Rectangle     { return w.img.Bounds() }
func (w *wrapper) At(x, y int) color.Color     { return w.img.At(x, y) }
func (w *wrapper) Set(x, y int, c color.Color) {}

// Immutable returns an immutable image
func Immutable(img image.Image) Image {
	return &wrapper{img: img}
}

// Crop will crop an image to fit the required bounds
func Crop(src image.Image, bounds image.Rectangle) Image {
	dstImage := image.NewRGBA(bounds.Sub(bounds.Min))
	draw.Draw(dstImage, dstImage.Bounds(), src, bounds.Min, draw.Src)
	return dstImage
}

// Expand will add the specified number of pixels to the edges of an Image
func Expand(src image.Image, top, left, bottom, right int) Image {
	// Ensure parameters are >=0
	left = util.Max(left, 0)
	right = util.Max(right, 0)
	top = util.Max(top, 0)
	bottom = util.Max(bottom, 0)

	// Calculate new image size
	oldBounds := src.Bounds()
	newBounds := image.Rectangle{Min: oldBounds.Min, Max: image.Point{X: oldBounds.Max.X + left + right, Y: oldBounds.Max.Y + top + bottom}}

	// Rectangle to draw old image into new one
	topLeft := image.Point{X: left, Y: top}
	drawBounds := oldBounds.Add(topLeft)

	dstImage := image.NewRGBA(newBounds)
	draw.Draw(dstImage, drawBounds, src, image.Point{}, draw.Src)

	return dstImage
}
