package util

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"github.com/peter-mount/go-script/packages"
	"image/color"
	"image/draw"
)

func init() {
	packages.RegisterPackage(&Util{})
}

type Util struct {
}

func (_ *Util) Rect(x1, y1, x2, y2 float64) util.Rectangle {
	return util.Rect(x1, y1, x2, y2)
}

func (_ *Util) GetStringBounds(gc *draw2dimg.GraphicContext, s string) util.Rectangle {
	return util.GetStringBounds(gc, s)
}

func (_ *Util) StringSize(gc *draw2dimg.GraphicContext, s string, a ...interface{}) util.Rectangle {
	return util.StringSize(gc, s, a...)
}

func (_ *Util) FitString(l, t, r, b, sl, st, sr, sb float64) (float64, float64, float64, float64) {
	return util.FitString(l, t, r, b, sl, st, sr, sb)
}

func (_ *Util) DrawStringLeft(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	return util.DrawStringLeft(gc, x, y, s, a...)
}

func (_ *Util) DrawStringCenter(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	return util.DrawStringCenter(gc, x, y, s, a...)
}

func (_ *Util) DrawStringRight(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	return util.DrawStringRight(gc, x, y, s, a...)
}

func (_ *Util) FloatToA(v float64) string {
	return util.FloatToA(v)
}

func (_ *Util) DrawColourBars(gc *draw2dimg.GraphicContext, bounds util.Rectangle, cols ...color.Color) (float64, float64) {
	return util.DrawColourBars(gc, bounds, cols...)
}

func (_ *Util) DrawColourBarsVertical(gc *draw2dimg.GraphicContext, bounds util.Rectangle, cols ...color.Color) (float64, float64) {
	return util.DrawColourBarsVertical(gc, bounds, cols...)
}

func (_ *Util) NewGraphicContext(img draw.Image) *draw2dimg.GraphicContext {
	return draw2dimg.NewGraphicContext(img)
}
