package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/util"
	"image"
	"image/color"
	"strings"
)

type Painter func(*draw2dimg.GraphicContext)

// Component is an entity within a frame
type Component interface {
	GetType() string
	Bounds() image.Rectangle
	SetBounds(image.Rectangle)
	Draw(draw2d.GraphicContext)
	Layout(draw2d.GraphicContext) bool
	IsUpdateRequired() bool
	Width() int
	Height() int
	Inset(int)
	Align(string)
	Font(string)
	SetType(string)
	Fill(color.Color)
	Stroke(color.Color)
	StrokeFill(color.Color)
	LineWidth(float64)
}

type BaseComponent struct {
	Type           string
	bounds         image.Rectangle // bounds of this container
	painter        Painter         // function to render this component
	font           string          // Font to use
	alignment      util.Alignment  // Text Alignment
	updateRequired bool
	fill           color.Color
	stroke         color.Color
	lineWidth      float64
	insetX         int // Inset on x-axis
	insetMinY      int // Inset on y-axis at the top, can differ to insetMaxY when dealing with titles
	insetMaxY      int // Inset on y-axis at the bottom
}

func (c *BaseComponent) SetPainter(painter Painter) {
	c.painter = painter
}

func (c *BaseComponent) IsUpdateRequired() bool {
	return c.updateRequired
}

func (c *BaseComponent) SetType(t string) {
	c.Type = t
}

func (c *BaseComponent) GetType() string {
	return c.Type
}

func (c *BaseComponent) Inset(inset int) {
	c.insetX = inset
	c.insetMinY = inset
	c.insetMaxY = inset
}

func (c *BaseComponent) Stroke(col color.Color) {
	c.stroke = col
	c.updateRequired = true
}

func (c *BaseComponent) StrokeFill(col color.Color) {
	c.stroke = col
	c.fill = col
	c.updateRequired = true
}

func (c *BaseComponent) Fill(col color.Color) {
	c.fill = col
	c.updateRequired = true
}

func (c *BaseComponent) LineWidth(f float64) {
	c.lineWidth = f
	c.updateRequired = true
}

func (c *BaseComponent) Font(font string) {
	c.font = font
	c.updateRequired = true
}

func (c *BaseComponent) Bounds() image.Rectangle {
	return c.bounds
}

func (c *BaseComponent) InsetBounds() image.Rectangle {
	return image.Rect(c.bounds.Min.X+c.insetX, c.bounds.Min.Y+c.insetMinY, c.bounds.Max.X-c.insetX, c.bounds.Max.Y-c.insetMaxY)
}

func (c *BaseComponent) GetInsets() (int, int, int, int) {
	return c.insetX, c.insetMinY, c.insetX, c.insetMaxY
}

func (c *BaseComponent) SetBounds(b image.Rectangle) {
	c.bounds = b
	c.updateRequired = true
	//fmt.Printf("SetBounds %q %v\n", c.Type, c.bounds)
}

func (c *BaseComponent) Width() int {
	return c.bounds.Dx()
}

func (c *BaseComponent) Height() int {
	return c.bounds.Dy()
}

// Layout the component. returns true if it has changed something
func (c *BaseComponent) Layout(_ draw2d.GraphicContext) bool {
	c.updateRequired = false
	return false
}

func (c *BaseComponent) Draw(ctx draw2d.GraphicContext) {
	c.paint(ctx.(*draw2dimg.GraphicContext), c.painter)
}

func (c *BaseComponent) paint(gc *draw2dimg.GraphicContext, painter Painter) {
	if c.painter != nil {
		gc.Save()
		defer gc.Restore()
		//fmt.Printf("translate %q\t%d, %d\n", c.Type, c.bounds.Min.X+c.inset, c.bounds.Min.Y+c.inset)
		gc.Translate(float64(c.bounds.Min.X+c.insetX), float64(c.bounds.Min.Y+c.insetMinY))

		if c.font != "" {
			_ = graph.SetFont(gc, c.font)
		}

		if c.stroke != nil {
			gc.SetStrokeColor(c.stroke)
		}

		if c.fill != nil {
			gc.SetFillColor(c.fill)
		}

		if c.lineWidth > 0.0 {
			gc.SetLineWidth(c.lineWidth)
		}

		painter(gc)
	}
}

func (c *BaseComponent) Align(s string) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "left":
		c.alignment = util.LeftAlignment
	case "right":
		c.alignment = util.RightAlignment
	case "center":
		c.alignment = util.CenterAlignment
	default:
		c.alignment = util.LeftAlignment
	}
}
