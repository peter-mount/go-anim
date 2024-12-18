package util

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"math"
	"strconv"
	"strings"
)

func GetStringBounds(gc *draw2dimg.GraphicContext, s string) Rectangle {
	sl, st, sr, sb := gc.GetStringBounds(s)
	return Rect(sl, st, sr, sb)
}

func StringSize(gc *draw2dimg.GraphicContext, s string, a ...interface{}) Rectangle {
	var rect Rectangle
	for i, str := range strings.Split(fmt.Sprintf(s, a...), "\n") {
		rect1 := GetStringBounds(gc, str)
		if i == 0 {
			rect = rect1
		} else {
			rect = rect.Add(rect1)
		}
	}
	return rect
}

func FitString(l, t, r, b, sl, st, sr, sb float64) (float64, float64, float64, float64) {
	return math.Min(l, sl), math.Min(t, st), math.Max(r, sr), math.Max(b, sb)
}

func DrawStringLeft(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	for _, str := range strings.Split(fmt.Sprintf(s, a...), "\n") {
		_, st, _, sb := gc.GetStringBounds(str)
		gc.FillStringAt(str, x, y+(sb-st)/2)
		y = y - st + sb
	}
	return y
}

func DrawStringCenter(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	for _, str := range strings.Split(fmt.Sprintf(s, a...), "\n") {
		sl, st, sr, sb := gc.GetStringBounds(str)
		gc.FillStringAt(str, x-(sr-sl)/2, y+(sb-st)/2)
		y = y - st + sb
	}
	return y
}

func DrawStringRight(gc *draw2dimg.GraphicContext, x, y float64, s string, a ...interface{}) float64 {
	for _, str := range strings.Split(fmt.Sprintf(s, a...), "\n") {
		sl, st, sr, sb := gc.GetStringBounds(str)
		gc.FillStringAt(str, x-(sr-sl), y+(sb-st)/2)
		y = y - st + sb
	}
	return y
}

func FloatToA(v float64) string {
	return strings.TrimSuffix(strconv.FormatFloat(v, 'f', 2, 64), ".00")
}
