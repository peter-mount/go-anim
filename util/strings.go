package util

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image"
	"image/color"
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

type Alignment uint8

const (
	LeftAlignment Alignment = iota
	CenterAlignment
	RightAlignment
)

// Fill fills the provided string based on this Alignment.
// If the string contains "\n" then it will be split and rendered as multiple lines.
//
// gc 			GraphicContext to draw to
// bounds 		image.Rectangle of the area to contain the string
// lineSpacing 	space to add between lines
// format,args	passed to fmt.Sprintf() before rendering
func (a Alignment) Fill(gc *draw2dimg.GraphicContext, bounds image.Rectangle, lineSpacing float64, format string, args ...interface{}) float64 {
	return a.Metrics(gc, bounds, lineSpacing, fmt.Sprintf(format, args...)).Fill(gc)
}

// Stroke the string based on this Alignment.
// If the string contains "\n" then it will be split and rendered as multiple lines.
//
// gc 			GraphicContext to draw to
// bounds 		image.Rectangle of the area to contain the string
// lineSpacing 	space to add between lines
// format,args	passed to fmt.Sprintf() before rendering
func (a Alignment) Stroke(gc *draw2dimg.GraphicContext, bounds image.Rectangle, lineSpacing float64, format string, args ...interface{}) float64 {
	return a.Metrics(gc, bounds, lineSpacing, fmt.Sprintf(format, args...)).Stroke(gc)
}

// FillStroke fills then strokes the string based on this Alignment.
// If the string contains "\n" then it will be split and rendered as multiple lines.
//
// gc 			GraphicContext to draw to
// bounds 		image.Rectangle of the area to contain the string
// lineSpacing 	space to add between lines
// format,args	passed to fmt.Sprintf() before rendering
func (a Alignment) FillStroke(gc *draw2dimg.GraphicContext, bounds image.Rectangle, lineSpacing float64, format string, args ...interface{}) float64 {
	return a.Metrics(gc, bounds, lineSpacing, fmt.Sprintf(format, args...)).FillStroke(gc)
}

type AlignmentMetrics struct {
	Bounds        image.Rectangle // Bounds of container
	ContentBounds image.Rectangle // Bounds of string within container
	MaxLineHeight float64         // Max line height over all lines
	MaxLineWidth  float64         // Max width of all lines
	BaseLines     []float64       // baseline for each line
	Widths        []float64       // Width of each line
	Lines         []string        // Line strings
	xFunc         func(int) float64
}

// Merge ensures that the line heights of both sets of metrics are the same.
// This is useful for when two sets need to line up with each other
func (m *AlignmentMetrics) Merge(b *AlignmentMetrics) {
	al := len(m.BaseLines)
	bl := len(b.BaseLines)
	ml := max(al, bl)

	var nb []float64

	for i := 0; i < ml; i++ {
		switch {
		case i < al && i < bl:
			nb = append(nb, max(m.BaseLines[i], b.BaseLines[i]))
		case i < al:
			nb = append(nb, m.BaseLines[i])
		case i < bl:
			nb = append(nb, b.BaseLines[i])
		}
	}

	m.BaseLines = nb
	b.BaseLines = nb

	m.MaxLineHeight = max(m.MaxLineHeight, b.MaxLineHeight)
	b.MaxLineHeight = m.MaxLineHeight
}

func (a Alignment) Metrics(gc *draw2dimg.GraphicContext, bounds image.Rectangle, lineSpacing float64, format string, args ...interface{}) *AlignmentMetrics {
	lineSpacing = max(0, lineSpacing)

	m := &AlignmentMetrics{
		Bounds: bounds,
		Lines:  strings.Split(fmt.Sprintf(format, args...), "\n"),
	}

	for _, line := range m.Lines {
		l, t, r, b := gc.GetStringBounds(line)
		m.MaxLineHeight = max(m.MaxLineHeight, b-t)
		m.MaxLineWidth = max(m.MaxLineWidth, r-l)
		m.BaseLines = append(m.BaseLines, -t) // t is always <=0
		m.Widths = append(m.Widths, r-l)
	}

	m.MaxLineHeight = m.MaxLineHeight + lineSpacing

	mlw := int(m.MaxLineWidth)
	lc := len(m.Lines)
	x0, y0, x1 := bounds.Min.X, bounds.Min.Y, bounds.Max.X
	cx := x0 + (bounds.Dx() >> 1)
	maxY := bounds.Min.Y + (lc * int(m.MaxLineHeight))
	switch a {
	case LeftAlignment:
		m.xFunc = m.leftX
		m.ContentBounds = image.Rect(x0, y0, x0+mlw, maxY)

	case CenterAlignment:
		m.xFunc = m.centerX
		lw := mlw >> 1
		m.ContentBounds = image.Rect(cx-lw, y0, cx+lw, maxY)

	case RightAlignment:
		m.xFunc = m.rightX
		m.ContentBounds = image.Rect(x1-mlw, y0, x1, maxY)

	default:
		// Should never happen but any invalid Alignment value will default to AlignLeft
		m.xFunc = m.leftX
		m.ContentBounds = image.Rect(x0, y0, x0+mlw, maxY)
	}

	return m
}

func (m *AlignmentMetrics) leftX(int) float64 { return float64(m.Bounds.Min.X) }

func (m *AlignmentMetrics) centerX(i int) float64 {
	return (float64(m.Bounds.Dx()) / 2) - (m.Widths[i] / 2)
}

func (m *AlignmentMetrics) rightX(i int) float64 { return float64(m.Bounds.Dx()) - m.Widths[i] }

// Fill fills the string defined in this AlignmentMetrics into the supplied GraphicContext
func (m *AlignmentMetrics) Fill(gc *draw2dimg.GraphicContext) float64 {
	gc.Save()
	gc.SetStrokeColor(color.White)
	gc.SetLineWidth(1.0)
	draw2dkit.Rectangle(gc, float64(m.ContentBounds.Min.X), float64(m.ContentBounds.Min.Y), float64(m.ContentBounds.Max.X), float64(m.ContentBounds.Max.Y))
	gc.Restore()
	return m.paint(gc.FillStringAt)
}

// Stroke the string defined in this AlignmentMetrics into the supplied GraphicContext
func (m *AlignmentMetrics) Stroke(gc *draw2dimg.GraphicContext) float64 {
	return m.paint(gc.StrokeStringAt)
}

// FillStroke first fills then strokes the string defined in this AlignmentMetrics into the supplied GraphicContext
func (m *AlignmentMetrics) FillStroke(gc *draw2dimg.GraphicContext) float64 {
	_ = m.Fill(gc)
	return m.Stroke(gc)
}

func (m *AlignmentMetrics) paint(f func(string, float64, float64) float64) float64 {
	y := 0.0

	for i, line := range m.Lines {
		f(line, m.xFunc(i), y+m.BaseLines[i])
		y = y + m.MaxLineHeight
	}

	return y
}

func FloatToA(v float64) string {
	return strings.TrimSuffix(strconv.FormatFloat(v, 'f', 2, 64), ".00")
}
