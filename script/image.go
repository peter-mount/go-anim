package script

import (
	"github.com/peter-mount/go-anim/renderer"
	"image"
	"image/color"
	"image/png"
	"io"
)

type Image struct {
	Width4K  int
	Height4K int
}

const (
	Width4K  = 3840 // 4K resolution
	Height4K = 2160 // 4K resolution
)

func newImage() *Image {
	return &Image{
		Width4K:  3840,
		Height4K: 2160,
	}
}

// New4K creates a new RGBA image at 4K resolution
func (_ Image) New4K() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width4K, Height4K))
}

// NewRGBA creates an RGBA image with the specified dimensions
func (_ Image) NewRGBA(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// Fill fills the image in the context with a specific colour
func (_ Image) Fill(ctx renderer.Context, background color.Color) {
	gc := ctx.Gc()
	gc.Save()
	defer gc.Restore()
	gc.SetFillColor(background)
	gc.ClearRect(0, 0, ctx.Width(), ctx.Height())
}

// WritePNG writes a PNG to a writer
func (i Image) WritePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
