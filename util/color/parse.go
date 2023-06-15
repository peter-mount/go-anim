package color

import (
	"encoding/xml"
	"errors"
	"fmt"
	"git.area51.dev/peter/videoident/util"
	"golang.org/x/image/colornames"
	"image/color"
	"strings"
)

// The MIT License (MIT) Copyright (c) 2019-2021 artipie.com
// based loosely on https://github.com/g4s8/hexcolor/blob/master/LICENSE

var ErrInvalidFormat = errors.New("invalid hex color format")

func Equals(a, b color.Color) bool {
	if a == nil {
		return b == nil
	}
	r0, g0, b0, a0 := a.RGBA()
	r1, g1, b1, a1 := b.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}

func AppendColorAttr(a []xml.Attr, n string, c color.Color) []xml.Attr {
	if c != nil {
		return util.AppendAttr(a, n, ColourString(c))
	}
	return a
}

func AppendColorSizeAttr(a []xml.Attr, n string, f float64, c color.Color) []xml.Attr {
	if f != 0 && c != nil {
		return util.AppendAttrf(a, n, "%s %s", util.FloatToA(f), ColourString(c))
	}
	return a
}

func ColourString(c color.Color) string {
	if c == nil {
		return "null"
	}

	r, g, b, a := c.RGBA()

	if a == 0 {
		return "transparent"
	}

	for n, cn := range colornames.Map {
		r1, g1, b1, a1 := cn.RGBA()

		if r == r1 && g == g1 && b == b1 && a == a1 {
			return n
		}
	}

	return fmt.Sprintf("#%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
}

func ParseColour(hex string) (c color.RGBA, err error) {
	if hex == "" {
		return c, ErrInvalidFormat
	}

	// If not a hex string then try to resolve it by name
	if hex[0] != '#' {
		name := strings.ToLower(strings.TrimSpace(hex))

		if name == "transparent" {
			return color.RGBA{A: 0}, nil
		}

		col, exists := colornames.Map[name]
		if !exists {
			return c, ErrInvalidFormat
		}
		return col, nil
	}

	c.A = 0xff

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return 10 + b - 'a'
		case b >= 'A' && b <= 'F':
			return 10 + b - 'A'
		}
		err = ErrInvalidFormat
		return 0
	}

	switch len(hex) {
	case 9:
		c.R = (hexToByte(hex[1]) << 4) + hexToByte(hex[2])
		c.G = (hexToByte(hex[3]) << 4) + hexToByte(hex[4])
		c.B = (hexToByte(hex[5]) << 4) + hexToByte(hex[6])
		c.A = (hexToByte(hex[7]) << 4) + hexToByte(hex[8])
	case 7:
		c.R = (hexToByte(hex[1]) << 4) + hexToByte(hex[2])
		c.G = (hexToByte(hex[3]) << 4) + hexToByte(hex[4])
		c.B = (hexToByte(hex[5]) << 4) + hexToByte(hex[6])
	case 5:
		c.R = hexToByte(hex[1]) * 17
		c.G = hexToByte(hex[2]) * 17
		c.B = hexToByte(hex[3]) * 17
		c.A = hexToByte(hex[4]) * 17
	case 4:
		c.R = hexToByte(hex[1]) * 17
		c.G = hexToByte(hex[2]) * 17
		c.B = hexToByte(hex[3]) * 17
	default:
		err = ErrInvalidFormat
	}
	return
}
