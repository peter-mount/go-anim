package util

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"github.com/peter-mount/go-script/packages"
	"image/color"
	"image/draw"
)

func init() {
	packages.RegisterPackage(&AnimUtil{})
}

type AnimUtil struct {
}

func (_ AnimUtil) Rect(x1, y1, x2, y2 float64) util.Rectangle {
	return util.Rect(x1, y1, x2, y2)
}

func (_ AnimUtil) GetStringBounds(gc *draw2dimg.GraphicContext, s string) util.Rectangle {
	return util.GetStringBounds(gc, s)
}

func (_ AnimUtil) StringSize(gc *draw2dimg.GraphicContext, s string, a ...interface{}) util.Rectangle {
	return util.StringSize(gc, s, a...)
}

func (_ AnimUtil) FitString(l, t, r, b, sl, st, sr, sb float64) (float64, float64, float64, float64) {
	return util.FitString(l, t, r, b, sl, st, sr, sb)
}

func (_ AnimUtil) DrawString(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	return util.DrawString(gc, x, y, s, a...)
}

func (_ AnimUtil) FloatToA(v float64) string {
	return util.FloatToA(v)
}

func (_ AnimUtil) DrawColourBars(gc *draw2dimg.GraphicContext, bounds util.Rectangle, cols ...color.Color) (float64, float64) {
	return util.DrawColourBars(gc, bounds, cols...)
}

func (_ AnimUtil) DrawColourBarsVertical(gc *draw2dimg.GraphicContext, bounds util.Rectangle, cols ...color.Color) (float64, float64) {
	return util.DrawColourBarsVertical(gc, bounds, cols...)
}

func (_ AnimUtil) NewGraphicContext(img draw.Image) *draw2dimg.GraphicContext {
	return draw2dimg.NewGraphicContext(img)
}
