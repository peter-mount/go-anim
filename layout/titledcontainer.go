package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/util"
	color2 "github.com/peter-mount/go-anim/util/color"
	"image"
	"image/color"
)

func TitledContainer(title string) Container {
	c := &titledContainer{
		rowContainer: *newRowContainer(),
		title:        title,
		titleFont:    "luxi 20 mono bold",
	}
	c.BaseComponent.Type = "TitledContainer"
	c.BaseComponent.painter = c.paint
	c.background, _ = color2.ParseColour("#0f2a3f")
	c.textColour, _ = color2.ParseColour("#7aa8ca")
	c.menuColor, _ = color2.ParseColour("#133b55")
	return c
}

type titledContainer struct {
	rowContainer
	title      string
	titleFont  string
	metrics    *util.AlignmentMetrics
	menuColor  color.RGBA
	textColour color.RGBA
	background color.RGBA
}

func (c *titledContainer) Title(s string) Container {
	c.title = s
	return c
}

func (c *titledContainer) TitleFont(font string) Container {
	c.titleFont = font
	return c
}

func (c *titledContainer) Layout(ctx draw2d.GraphicContext) bool {
	c.BaseComponent.paint(ctx.(*draw2dimg.GraphicContext), func(gc *draw2dimg.GraphicContext) {
		gc.Save()
		_ = graph.SetFont(gc, c.titleFont)
		b := c.Bounds()
		b = b.Add(image.Pt(-b.Min.X, -b.Min.Y))
		c.metrics = util.LeftAlignment.Metrics(gc, b, 2, c.title)
		gc.Restore()

		c.insetMinY = 10 + int(c.metrics.MaxLineHeight)
		c.insetMaxY = 10
		c.insetX = 5

		_ = c.rowContainer.Layout(ctx)

		c.SetBounds(c.rowContainer.BaseComponent.Bounds())
	})
	return true
}

func (c *titledContainer) paint(gc *draw2dimg.GraphicContext) {
	gc.Save()
	b := c.Bounds()

	gc.SetFillColor(c.background)
	gc.BeginPath()
	draw2dkit.Rectangle(gc,
		float64(c.insetX), float64(0),
		float64(b.Dx()-c.insetX-c.insetX), float64(b.Dy()-c.insetMaxY),
	)
	gc.Fill()

	gc.SetFillColor(c.menuColor)
	gc.BeginPath()
	draw2dkit.Rectangle(gc,
		float64(c.insetX), float64(0),
		float64(b.Dx()-c.insetX-c.insetX), float64(c.insetMinY),
	)
	gc.Fill()

	gc.SetFillColor(c.textColour)
	_ = graph.SetFont(gc, c.titleFont)
	gc.Translate(float64(c.insetX)+5, 5)
	c.metrics.Fill(gc)

	gc.Restore()

	gc.SetFillColor(c.textColour)

	c.rowContainer.paint(gc)

	gc.Save()
	gc.SetStrokeColor(c.background)
	gc.SetLineWidth(2)
	gc.BeginPath()
	draw2dkit.Rectangle(gc,
		float64(c.insetX), float64(0),
		float64(b.Dx()-c.insetX-c.insetX), float64(b.Dy()-c.insetMaxY),
	)
	gc.Stroke()
	gc.Restore()

}
