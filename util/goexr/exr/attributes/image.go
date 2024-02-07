package attributes

import (
	"image"
	"image/color"
	"image/draw"
)

type Image interface {
	ImageAttributes
	draw.Image
}

// Wrap any draw.Image and return one which supports Attributes
func Wrap(img draw.Image) Image {
	if i, ok := img.(Image); ok {
		return i
	}
	return &wrapper{img: img}
}

type wrapper struct {
	DefaultImageAttributes
	img draw.Image
}

func (i *wrapper) ColorModel() color.Model {
	return i.img.ColorModel()
}

func (i *wrapper) Bounds() image.Rectangle {
	return i.img.Bounds()
}

func (i *wrapper) At(x, y int) color.Color {
	return i.img.At(x, y)
}

func (i *wrapper) Set(x, y int, c color.Color) {
	i.img.Set(x, y, c)
}
