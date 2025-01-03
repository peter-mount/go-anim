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

// AutoCrop will crop an image removing anything that is black around the image.
func AutoCrop(src image.Image) Image {
	b := src.Bounds()

	// Look for first line
	for y := b.Min.Y; y <= b.Max.Y; y++ {
		if !isBlackRow(src, y, b.Min.X, b.Max.X) {
			b.Min.Y = y
			break
		}
	}

	// Look for last line
	for y := b.Max.Y; y >= b.Min.Y; y-- {
		if !isBlackRow(src, y, b.Min.X, b.Max.X) {
			b.Max.Y = y
			break
		}
	}

	// First column
	for x := b.Min.X; x <= b.Max.X; x++ {
		if !isBlackCol(src, x, b.Min.Y, b.Max.Y) {
			b.Min.X = x
			break
		}
	}

	// Last column
	for x := b.Max.X; x > b.Min.X; x-- {
		if !isBlackCol(src, x, b.Min.Y, b.Max.Y) {
			b.Max.X = x
			break
		}
	}

	if b.Dx() > 0 && b.Dy() > 0 {
		return Crop(src, b)
	}

	// Nothing to crop so return a duplicate but writable image
	return DuplicateImage(src)
}

func IsBlack(c color.Color) bool {
	return !IsNotBlack(c)
}

func IsNotBlack(c color.Color) bool {
	r1, g1, b1, a1 := c.RGBA()
	return a1 == 0 || (r1 != 0 && g1 != 0 && b1 != 0)
}

func isBlackRow(src image.Image, y, x0, x1 int) bool {
	for x := x0; x <= x1; x++ {
		if IsNotBlack(src.At(x, y)) {
			return false
		}
	}
	return true
}

func isBlackCol(src image.Image, x, y0, y1 int) bool {
	for y := y0; y <= y1; y++ {
		if IsNotBlack(src.At(x, y)) {
			return false
		}
	}
	return true
}

func Mask(img, mask image.Image) (Image, error) {
	return Of(func(x, y int, col color.Color) (color.Color, error) {
		if IsNotBlack(mask.At(x, y)) {
			return col, nil
		}
		return color.Black, nil
	}).
		DoNew(img)
}

// DrawMask draws img over dest using mask to select which pixels in img are to be
// copied over. This returns a new image
func DrawMask(img, mask, dest image.Image) (Image, error) {
	return Of(func(x, y int, col color.Color) (color.Color, error) {
		if IsNotBlack(mask.At(x, y)) {
			return img.At(x, y), nil
		}
		return col, nil
	}).
		DoNew(dest)
}
