package image

import (
	"github.com/peter-mount/go-anim/graph"
	color2 "github.com/peter-mount/go-anim/graph/color"
	"github.com/peter-mount/go-anim/graph/filter"
	"github.com/peter-mount/go-anim/renderer"
	"github.com/peter-mount/go-anim/script/render"
	"github.com/peter-mount/go-anim/util/goexr/exr"
	"github.com/peter-mount/go-script/packages"
	"image"
	"image/color"
)

func init() {
	p := &Image{
		Width4K:         3840,
		Height4K:        2160,
		Width1080p:      Width1080p,
		Height1080p:     Height1080p,
		Width720p:       Width720p,
		Height720p:      Height720p,
		BlackMaskFilter: filter.BlackMaskFilter,
		WhiteMaskFilter: filter.WhiteMaskFilter,

		encoders: map[string]render.Encoder{
			".png":  &render.PNG{},
			".jpg":  &render.JPEG{},
			".jpeg": &render.JPEG{},
			".tiff": &render.TIFF{},
			".tif":  &render.TIFF{},
		},

		decoders: map[string]render.Decoder{
			".png":  &render.PNG{},
			".jpg":  &render.JPEG{},
			".jpeg": &render.JPEG{},
			".tiff": &render.TIFF{},
			".tif":  &render.TIFF{},
		},
	}

	packages.RegisterPackage(p)
}

type Image struct {
	Width4K         int
	Height4K        int
	Width1080p      int
	Height1080p     int
	Width720p       int
	Height720p      int
	BlackMaskFilter graph.Filter
	WhiteMaskFilter graph.Filter

	encoders map[string]render.Encoder
	decoders map[string]render.Decoder
}

const (
	Width4K     = 3840 // 4K resolution, 2160p
	Height4K    = 2160 // 4K resolution, 2160p
	Width1080p  = 1920 // FHD 1080p resolution
	Height1080p = 1080 // FHD 1080p resolution
	Width720p   = 1280 // HD 720p resolution
	Height720p  = 720
)

// New returns a new RGBA image with the specified dimensions
func (_ *Image) New(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// New4K creates a new RGBA image at 4K resolution
func (_ *Image) New4K() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width4K, Height4K))
}

// New2160p creates a new RGBA image at 2160p.
// This is the same as New4K.
func (g *Image) New2160p() *image.RGBA {
	return g.New4K()
}

// New1080p creates a new RGBA image at 1080p resolution,
// also known as FHD or Full HD.
func (_ *Image) New1080p() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width1080p, Height1080p))
}

// New720p creates a new RGBA image at 720p resolution,
// also known has HD or HD Ready.
func (_ *Image) New720p() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, Width720p, Height720p))
}

// NewRGBA creates an RGBA image with the specified dimensions
func (_ *Image) NewRGBA(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// NewRGBA64 creates an RGBA64 image with the specified dimensions
func (_ *Image) NewRGBA64(w, h int) *image.RGBA64 {
	return image.NewRGBA64(image.Rect(0, 0, w, h))
}

// NewFloat16 returns a RGBAImage using float16 for each colour component
func (_ *Image) NewFloat16(w, h int) *exr.RGBAImage {
	return exr.NewFloat16(image.Rect(0, 0, w, h))
}

// NewFloat32 returns a RGBAImage using float32 for each colour component
func (_ *Image) NewFloat32(w, h int) *exr.RGBAImage {
	return exr.NewFloat32(image.Rect(0, 0, w, h))
}

// Fill fills the image in the context with a specific colour
func (_ *Image) Fill(ctx renderer.Context, background color.Color) {
	gc := ctx.Gc()
	gc.Save()
	defer gc.Restore()
	gc.SetFillColor(background)
	gc.ClearRect(0, 0, ctx.Width(), ctx.Height())
}

func (_ *Image) Histogram(src image.Image) *color2.Histogram {
	return color2.NewHistogram().AnalyzeImage(src)
}

func (_ *Image) Equalize(h *color2.Histogram, b image.Rectangle) graph.Filter {
	return filter.EqualizeFilter(h, b)
}

// Filter applies a graph.Filter on a source image within the specified bounds,
// writing the result to the destination image.
//
// The source and destination image may be the same Image if the filter supports it.
func (g *Image) Filter(f graph.Filter, src image.Image, dst graph.Image, b image.Rectangle) error {
	return f.Do(src, dst, b)
}

// FilterNew applies a graph.Filter on a source image,
// returning a new mutable image with the new content.
func (g *Image) FilterNew(f graph.Filter, src image.Image) (graph.Image, error) {
	return f.DoNew(src)
}

// FilterOver applies the filter over the supplied mutable image,
// overwriting its previous state.
func (g *Image) FilterOver(f graph.Filter, src graph.Image) error {
	return f.DoOver(src)
}

func (_ *Image) Duplicate(img image.Image) graph.Image {
	return graph.DuplicateImage(img)
}

func (_ *Image) Immutable(img image.Image) graph.Image {
	return graph.Immutable(img)
}
