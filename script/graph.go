package script

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/renderer"
	draw2d2 "github.com/peter-mount/go-anim/util/draw2d"
	"github.com/peter-mount/go-anim/util/font"
	"image"
	"image/color"
)

type Graph struct {
}

func (_ Graph) NewContext() renderer.Context {
	return renderer.NewContext()
}

func (_ Graph) NewSizedContext(w, h int) renderer.Context {
	ctx := renderer.NewContext()
	ctx.SetImage(image.NewRGBA(image.Rect(0, 0, w, h)))
	return ctx
}

func (_ Graph) NewFont(name string, size float64, family draw2d.FontFamily, style draw2d.FontStyle) font.Font {
	return font.New(name, size, family, style)
}

func (_ Graph) ParseFont(s string) (font.Font, error) {
	return font.ParseFont(s)
}

func (_ Graph) SetFont(gc *draw2dimg.GraphicContext, s string) error {
	f, err := font.ParseFont(s)
	if err != nil {
		return err
	}

	gc.SetFontData(f.FontData())
	gc.SetFontSize(f.Size())
	return nil
}

func (_ Graph) FillPoly(gc *draw2dimg.GraphicContext, c color.Color, v ...float64) {
	draw2d2.FillPoly(gc, c, v...)
}

func (_ Graph) FillPolyRel(gc *draw2dimg.GraphicContext, c color.Color, v ...float64) {
	draw2d2.FillPolyRel(gc, c, v...)
}

func (_ Graph) FillRectangle(gc *draw2dimg.GraphicContext, x, y, w, h float64, c color.Color) (float64, float64) {
	return draw2d2.FillRectangle(gc, x, y, w, h, c)
}

func (_ Graph) Rectangle(gc *draw2dimg.GraphicContext, x, y, w, h float64) {
	draw2d2.Rectangle(gc, x, y, w, h)
}

func (_ Graph) RelLine(gc *draw2dimg.GraphicContext, x, y float64, v ...float64) {
	draw2d2.RelLine(gc, x, y, v...)
}
