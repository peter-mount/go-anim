package renderer

import (
	"github.com/golang/freetype/raster"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/util"
	"github.com/peter-mount/go-anim/util/goexr/exr"
	"golang.org/x/image/draw"
	"image"
)

// NewContext creates a new context.
// start,end are the frame range to cover, frameRate the frame rate of the track
// whilst duration is the total duration of the track in frames
func NewContext(width, height int) Context {
	return New(image.NewRGBA(image.Rect(0, 0, width, height)))
}

func New(img draw.Image) Context {
	ctx := &context{}
	return ctx.SetImage(img)
}

type context struct {
	img      draw.Image                // Image to use for frame generation
	width    int                       // Width of image
	height   int                       // Height of image
	gc       *draw2dimg.GraphicContext // Graphic context
	userdata map[string]any            // User data
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
	return c.SetImage(image.NewRGBA(image.Rect(0, 0, c.width-1, c.height-1)))
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
	return util.Rect(0, 0, float64(c.width-1), float64(c.height-1))
}

func (c *context) Gc() *draw2dimg.GraphicContext {
	return c.gc
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
	c.Gc().Restore()
	return nil
}

func (c *context) Reset() Context {
	// Reset the Context state
	var painter draw2dimg.Painter
	switch selectImage := c.img.(type) {

	case *image.RGBA:
		painter = raster.NewRGBAPainter(selectImage)

	case *exr.RGBAImage:
		painter = exr.NewRGBAImagePainter(selectImage)

	default:
		panic("Image type not supported")
	}

	c.gc = draw2dimg.NewGraphicContextWithPainter(c.img, painter)
	return c
}

func (c *context) Draw(d Drawable) {
	if d != nil {
		d.Draw(c)
	}
}

func (c *context) Filter(f graph.Filter) error {
	return f.DoOver(c.Image())
}

func (c *context) FilterBounds(f graph.Filter, b image.Rectangle) error {
	return f.Do(c.Image(), c.Image(), b)
}

func (c *context) Map(m graph.Mapper) error {
	return c.Filter(m.Filter())
}

func (c *context) MapBounds(m graph.Mapper, b image.Rectangle) error {
	return c.FilterBounds(m.Filter(), b)
}
