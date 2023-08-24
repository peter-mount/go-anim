package script

import (
	color2 "github.com/peter-mount/go-anim/graph/color"
	"github.com/peter-mount/go-anim/renderer"
	"image"
	"image/color"
)

type Image struct {
	Width4K  int
	Height4K int
}

const (
	Width4K     = 3840 // 4K resolution, 2160p
	Height4K    = 2160 // 4K resolution, 2160p
	Width1080p  = 1920 // FHD 1080p resolution
	Height1080p = 1080 // FHD 1080p resolution
	Width720p   = 1280 // HD 720p resolution
	Height720p  = 720
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

// New2160p creates a new RGBA image at 2160p.
// This is the same as New4K.
func (g Image) New2160p() *image.RGBA {
	return g.New4K()
}

// New1080p creates a new RGBA image at 1080p resolution,
// also known as FHD or Full HD.
func (_ Image) New1080p() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width1080p, Height1080p))
}

// New720p creates a new RGBA image at 720p resolution,
// also known has HD or HD Ready.
func (_ Image) New720p() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width720p, Height720p))
}

// NewRGBA creates an RGBA image with the specified dimensions
func (_ Image) NewRGBA(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewRGBA64 creates an RGBA64 image with the specified dimensions
func (_ Image) NewRGBA64(w, h int) *image.RGBA64 {
	return image.NewRGBA64(image.Rect(0, 0, w, h))
}

// Fill fills the image in the context with a specific colour
func (_ Image) Fill(ctx renderer.Context, background color.Color) {
	gc := ctx.Gc()
	gc.Save()
	defer gc.Restore()
	gc.SetFillColor(background)
	gc.ClearRect(0, 0, ctx.Width(), ctx.Height())
}

func (_ Image) Histogram(src image.Image) *color2.Histogram {
	return color2.NewHistogram().AnalyzeImage(src)
}
