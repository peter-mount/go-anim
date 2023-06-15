package renderer

import (
	"git.area51.dev/peter/videoident/util"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/draw"
	"image"
)

// NewContext creates a new context.
// start,end are the frame range to cover, frameRate the frame rate of the track
// whilst duration is the total duration of the track in frames
func NewContext(start, end int, frameRate, duration float64) Context {
	ctx := &context{
		start:     start,
		end:       end,
		nextFrame: start,
		frameRate: frameRate,
		duration:  duration,
	}
	return ctx.NewImage()
}

type context struct {
	curFrame  int                       // current frame number
	nextFrame int                       // next frame number
	start     int                       // start frame number
	end       int                       // end frame number
	frameRate float64                   // frame rate fps
	duration  float64                   // Duration of clip
	img       draw.Image                // Image to use for frame generation
	width     int                       // Width of image
	height    int                       // Height of image
	gc        *draw2dimg.GraphicContext // Graphic context
	userdata  map[string]any            // User data
}

func (c *context) HasNext() bool {
	return c.nextFrame <= c.end
}

func (c *context) Image() draw.Image {
	return c.img
}

func (c *context) SetImage(img draw.Image) Context {
	c.img = img
	b := img.Bounds()
	c.width = b.Dx()
	c.height = b.Dy()
	return c.Reset()
}

func (c *context) NewImage() Context {
	if c.width == 0 {
		c.width = util.Width4K
	}
	if c.height == 0 {
		c.height = util.Height4K
	}
	return c.SetImage(image.NewRGBA(image.Rect(0, 0, c.width, c.height)))
}

func (c *context) Width() int {
	return c.width
}

func (c *context) Height() int {
	return c.height
}

func (c *context) Center() (float64, float64) {
	return float64(c.width) / 2, float64(c.height) / 2
}

func (c *context) Bounds() util.Rectangle {
	return util.Rect(0, 0, float64(c.width), float64(c.height))
}

func (c *context) Gc() *draw2dimg.GraphicContext {
	return c.gc
}

func (c *context) FrameRate() float64 {
	return c.frameRate
}

func (c *context) Duration() float64 {
	return c.duration
}

func (c *context) Start() int {
	return c.start
}

func (c *context) End() int {
	return c.end
}

func (c *context) Frame() int {
	return c.curFrame
}

func (c *context) FrameF() float64 {
	return float64(c.Frame())
}

func (c *context) Get(k string) any {
	return c.userdata[k]
}

func (c *context) Set(k string, v any) Context {
	if v == nil {
		delete(c.userdata, k)
	} else {
		c.userdata[k] = v
	}
	return c
}

func (c *context) Remove(k string) Context {
	delete(c.userdata, k)
	return c
}

// Create from CreateCloser interface, used in try resources block
// to save and close state in the context
func (c *context) Create() error {
	c.Gc().Save()
	return nil
}

// Close from CreateCloser interface, used in try resources block
// to save and close state in the context
func (c *context) Close() error {
	c.Gc().Close()
	return nil
}

func (c *context) Reset() Context {
	// Reset the Context state
	c.gc = draw2dimg.NewGraphicContext(c.img)
	return c
}
