package renderer

import (
	"git.area51.dev/peter/videoident/util"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/draw"
)

type Context interface {
	// Img The underlying image
	Img() draw.Image
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
	// FrameRate of the video
	FrameRate() float64
	// Duration of the video
	Duration() float64
	// Start frame in this run
	Start() int
	// End frame in this run
	End() int
	// Frame Current frame number being rendered
	Frame() int
	// FrameF shorthand for float64(Frame())
	FrameF() float64
	// Get returns a named user object, used in keeping state.
	// This is cleared at the start of each frame
	Get(string) any
	// Set allows for a user object to be stored for retrieval with Get().
	// This allows for storing information during a frame's rendering.
	// This is cleared at the start of each frame
	Set(string, any) Context
	// Remove removes a key from the user object storage
	Remove(k string) Context
	// Render will call a Renderer with a clean instance of the context
	Render(Renderer) error
	// HasNext true if there's more frames in the Context
	HasNext() bool
}
