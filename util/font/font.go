package font

import (
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/math/fixed"
	"math"
	"strconv"
	"strings"
)

type Font interface {
	Name() string
	Size() float64
	Family() draw2d.FontFamily
	Style() draw2d.FontStyle
	FontData() draw2d.FontData
	Extents() draw2dimg.FontExtents
	// WithSize returns this Font but with the specified size
	WithSize(float64) Font
	// WithStyle returns this Font but with the specified style
	WithStyle(draw2d.FontStyle) Font
	// WithFamily returns this F|ont but with the specified family
	WithFamily(draw2d.FontFamily) Font
	// Set the font in a GraphicContext
	Set(*draw2dimg.GraphicContext)
	StringLength(s string) float64
}

type font struct {
	name     string
	size     float64
	family   draw2d.FontFamily
	style    draw2d.FontStyle
	fontData draw2d.FontData
}

func New(name string, size float64, family draw2d.FontFamily, style draw2d.FontStyle) Font {
	return &font{
		name:     name,
		size:     size,
		family:   family,
		style:    style,
		fontData: draw2d.FontData{Name: name, Family: family, Style: style},
	}
}

func (f *font) Name() string {
	return f.name
}

func (f *font) Size() float64 {
	return f.size
}

func (f *font) Family() draw2d.FontFamily {
	return f.family
}

func (f *font) Style() draw2d.FontStyle {
	return f.style
}

func (f *font) FontData() draw2d.FontData {
	return f.fontData
}

func (f *font) Extents() draw2dimg.FontExtents {
	return draw2dimg.Extents(draw2d.GetFont(f.fontData), f.size)
}

func (f *font) WithSize(size float64) Font {
	return New(f.name, size, f.family, f.style)
}

func (f *font) WithStyle(style draw2d.FontStyle) Font {
	return New(f.name, f.size, f.family, style)
}

func (f *font) WithFamily(family draw2d.FontFamily) Font {
	return New(f.name, f.size, family, f.style)
}

func (f *font) Set(gc *draw2dimg.GraphicContext) {
	gc.SetFontData(f.fontData)
	gc.SetFontSize(f.size)
}

func (f *font) StringLength(s string) float64 {
	tf := draw2d.GetFont(f.fontData)
	x := 0.0

	// Size is font size * 92DPI * (64/72)
	size := f.size * 92.0 * (64.0 / 72.0)

	prev, hasPrev := truetype.Index(0), false
	for _, fRune := range s {
		index := tf.Index(fRune)
		if hasPrev {
			x += FUnitsToFloat64(tf.Kern(fixed.Int26_6(size), prev, index))
		}
		x += FUnitsToFloat64(tf.HMetric(fixed.Int26_6(size), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x
}

func FUnitsToFloat64(x fixed.Int26_6) float64 {
	scaled := x << 2
	return float64(scaled/256) + float64(scaled%256)/256.0
}

func ParseFont(s string) (Font, error) {
	name := ""
	size := 10.0
	family := draw2d.FontFamilyMono
	style := draw2d.FontStyleNormal

	// TODO add some form of error handling here
	for _, e := range strings.Split(s, " ") {
		if vf, err := strconv.ParseFloat(e, 64); err == nil {
			// Set the size, enforce a minimum of 1 pixel
			size = math.Max(1, vf)
		} else {
			switch e {
			case "bold":
				style = style | draw2d.FontStyleBold
			case "italic":
				style = style | draw2d.FontStyleItalic
			case "sans":
				family = draw2d.FontFamilySans
			case "serif":
				family = draw2d.FontFamilySerif
			case "mono":
				family = draw2d.FontFamilyMono
			default:

				name = e
			}
		}
	}

	f := New(name, size, family, style)

	// This tests the font exists, returns an error if it doesn't or is corrupted
	_, err := draw2d.GetGlobalFontCache().Load(f.FontData())
	if err != nil {
		return nil, err
	}

	return f, nil
}
