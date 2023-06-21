package renderer

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/util"
	"golang.org/x/image/draw"
	"image"
)

type Context interface {
	Image() draw.Image
	SetImage(draw.Image) Context
	NewImage() Context
	// Width of the image
	Width() int
	// Height of the image
	Height() int
	// Bounds of the image
	Bounds() util.Rectangle
	// Center coordinates of the Image
	Center() (float64, float64)
	// Gc draw2dimg.GraphicContext
	Gc() *draw2dimg.GraphicContext
	// Get returns a named user object, used in keeping state.
	// This is cleared at the start of each frame
	Get(string) any
	// Set allows for a user object to be stored for retrieval with Get().
	// This allows for storing information during a frame's rendering.
	// This is cleared at the start of each frame
	Set(string, any) Context
	// Remove removes a key from the user object storage
	Remove(k string) Context
	Create() error
	Close() error
	Reset() Context
	Draw(d Drawable)
	Filter(f graph.Filter) error
	FilterBounds(f graph.Filter, b image.Rectangle) error
	Map(m graph.Mapper) error
	MapBounds(m graph.Mapper, b image.Rectangle) error
}

type Drawable interface {
	Draw(ctx Context)
}
