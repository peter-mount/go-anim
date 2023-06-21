// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package units defines units of length such as inches or pixels.
//
// Functions like Inches and Pixels return a Value in the corresponding unit.
// For example:
//
//	v := unit.Inches(4.5)
//
// represents four and a half inches.
//
// Converting between pixels (px), physical units (dp, pt, in, mm) and
// font-face-relative measures (em, ex, ch) depends on the context, such as the
// screen's DPI resolution and the active font face. That context is
// represented by the Converter type.
//
// Conversions may be lossy. Converting 4.5 inches to pixels and back may
// result in something slightly different than 4.5. Similarly, converting 4
// inches and 0.5 inches to pixels and then adding the results won't
// necessarily equal the conversion of 4.5 inches to pixels.
//
// Note that what CSS (Cascading Style Sheets) calls "px" differs from what
// this package calls "px". For legacy reasons, the CSS semantics are that 1
// inch should roughly equal 96csspx regardless of the actual DPI resolution,
// as per https://developer.mozilla.org/en/docs/Web/CSS/length. This package's
// semantics are that 1px means exactly one physical pixel, always. This
// package represents 1csspx as 1.666666667dp, since there are 160 density
// independent pixels per inch, the same definition as Android.
package unit

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util"
	"github.com/peter-mount/go-anim/util/font"
	"golang.org/x/image/math/fixed"
	"math"
	"strconv"
	"strings"
)

const (
	DensityIndependentPixelsPerInch = 160
	MillimetresPerInch              = 25.4
	PointsPerInch                   = 72.0
)

// Value is a number and a unit.
type Value struct {
	F float64
	U Unit
}

func (v Value) Convert(gc *draw2dimg.GraphicContext, to Unit) Value {
	if v.U == to {
		return v
	}

	return Value{
		F: v.F * to.PixelsPer(gc) / v.U.PixelsPer(gc),
		U: to,
	}
}

// Pixels returns the value in pixels as a float64
func (v Value) Pixels(gc *draw2dimg.GraphicContext) float64 {
	return v.Convert(gc, Px).F
}

func Pixels(f float64) Value {
	return Value{F: f, U: Px}
}

// PixelsPer returns the number of pixels in the unit u.
func (u Unit) PixelsPer(gc *draw2dimg.GraphicContext) float64 {
	dpi := 92.0
	if gc != nil {
		dpi = float64(gc.DPI)
	}

	switch u {
	case Px:
		return 1
	case Dp:
		return dpi / DensityIndependentPixelsPerInch
	case Pt:
		return dpi / PointsPerInch
	case Mm:
		return dpi / MillimetresPerInch
	case In:
		return dpi
	}

	if gc != nil {
		f := draw2d.GetFont(gc.GetFontData())
		ext := draw2dimg.Extents(f, gc.GetFontSize())
		h := ext.Height

		switch u {
		case Em:
			return h
		case Ex:
			return h / 2
		case Ch:
			if i := f.Index('0'); i > 0 {
				return font.FUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.GetFontSize()), i).AdvanceWidth)
			}
			return h / 2
		}
	}

	return 1
}

func (v Value) IsZero() bool {
	return v.F == 0
}

func (v Value) Equals(b Value) bool {
	return v.F == b.F && v.U == b.U
}

// String implements the fmt.Stringer interface.
func (v Value) String() string {
	return util.FloatToA(v.F) + names[v.U]
}

func ParseValue(s string) (Value, error) {
	for u, n := range names {
		if strings.HasSuffix(s, n) {
			f, err := strconv.ParseFloat(strings.TrimSuffix(s, n), 64)
			if err != nil {
				return Value{}, err
			}
			return Value{F: f, U: Unit(u)}, nil
		}
	}

	// Default to Px
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{F: f, U: Px}, nil
}

// Add adds two values. The result will be in the unit of the left hand side.
func (v Value) Add(gc *draw2dimg.GraphicContext, b Value) Value {
	return Value{
		F: v.F + b.Convert(gc, v.U).F,
		U: v.U,
	}
}

// Sub subtracts two values. The result will be in the unit of the left hand side.
func (v Value) Sub(gc *draw2dimg.GraphicContext, b Value) Value {
	return Value{
		F: v.F - b.Convert(gc, v.U).F,
		U: v.U,
	}
}

func (v Value) Min(gc *draw2dimg.GraphicContext, b Value) Value {
	return Value{
		F: math.Min(v.F, b.Convert(gc, v.U).F),
		U: v.U,
	}
}

func (v Value) Max(gc *draw2dimg.GraphicContext, b Value) Value {
	return Value{
		F: math.Max(v.F, b.Convert(gc, v.U).F),
		U: v.U,
	}
}
