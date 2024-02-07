package exr

import (
	"github.com/peter-mount/go-anim/util/goexr/exr/attributes"
	"image"
	"image/color"

	"github.com/peter-mount/go-anim/util/goexr/exr/internal/exr"
)

// RGBAImage represents an EXR image that consists of R, G, B, and A components.
//
// Even if the original image that is loaded does not contain all of the
// components, default ones will be assigned.
type RGBAImage struct {
	attributes.DefaultImageAttributes
	rect     image.Rectangle
	channelR exr.PixelData
	channelG exr.PixelData
	channelB exr.PixelData
	channelA exr.PixelData
}

func NewFloat32(rect image.Rectangle) *RGBAImage {
	return newRGBA(rect, exr.NewFloat32PixelData)
}

func NewFloat16(rect image.Rectangle) *RGBAImage {
	return newRGBA(rect, exr.NewFloat16PixelData)
}

func newRGBA(rect image.Rectangle, f func(window exr.Box2i, xSampling, ySampling int32) exr.PixelData) *RGBAImage {
	window := exr.Box2iFromRect(rect)
	i := &RGBAImage{
		rect:     rect,
		channelR: f(window, 1, 1),
		channelG: f(window, 1, 1),
		channelB: f(window, 1, 1),
		channelA: f(window, 1, 1),
	}

	return i
}

// ColorModel returns the RGBAImage's color model.
func (i *RGBAImage) ColorModel() color.Model {
	return RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i *RGBAImage) Bounds() image.Rectangle {
	return i.rect
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
//
// The returned color is of type RGBAColor which can be used to acquire the
// linear (float) components of the color.
func (i *RGBAImage) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}.In(i.rect)) {
		return RGBAColor{}
	}
	return RGBAColor{
		R: i.channelR.Float32(x, y),
		G: i.channelG.Float32(x, y),
		B: i.channelB.Float32(x, y),
		A: i.channelA.Float32(x, y),
	}
}

func (i *RGBAImage) GetRGBA(x, y int) (float32, float32, float32, float32) {
	if !(image.Point{X: x, Y: y}.In(i.rect)) {
		return 0, 0, 0, 0
	}
	return i.channelR.Float32(x, y),
		i.channelG.Float32(x, y),
		i.channelB.Float32(x, y),
		i.channelA.Float32(x, y)
}

func (i *RGBAImage) Set(x, y int, c color.Color) {
	if (image.Point{X: x, Y: y}.In(i.rect)) {
		// Convert to our colour model, which we know will always return RGBAColor
		rc := rgbaModel(c).(RGBAColor)
		i.channelR.Set(x, y, rc.R)
		i.channelG.Set(x, y, rc.G)
		i.channelB.Set(x, y, rc.B)
		i.channelA.Set(x, y, rc.A)
	}
}

func (i *RGBAImage) SetRGBA(x, y int, r, g, b, a float32) {
	if (image.Point{X: x, Y: y}.In(i.rect)) {
		i.channelR.Set(x, y, r)
		i.channelG.Set(x, y, g)
		i.channelB.Set(x, y, b)
		i.channelA.Set(x, y, a)
	}
}
