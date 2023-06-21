package color

import "image/color"

// Mono is a Mapper which converts an image into greyscale.
//
// It does this by just keeping the Green channel, as the human eye is more
// receptive to this channel.
//
// It's also how most colour scanners do greyscale by scanning in colour & taking
// just the green component.
func Mono(col color.Color) (color.Color, error) {
	return Green(col)
}

// Mono16 is a Mapper which converts an image into greyscale.
// The default implementationm is to use
// It does this by just keeping the Green channel, as the human eye is more
// receptive to this channel.
//
// It's also how most colour scanners do greyscale by scanning in colour & taking
// just the green component.
func Mono16(col color.Color) (color.Color, error) {
	return Green16(col)
}

func red(col color.Color) uint32 {
	v, _, _, _ := col.RGBA()
	return v
}

func green(col color.Color) uint32 {
	_, v, _, _ := col.RGBA()
	return v
}

func blue(col color.Color) uint32 {
	_, _, v, _ := col.RGBA()
	return v
}

// Red maps to just the red component in 8 bits
func Red(col color.Color) (color.Color, error) {
	return color.Gray{Y: uint8(red(col) >> 8)}, nil
}

// Red16 maps to just the red component in 16 bits
func Red16(col color.Color) (color.Color, error) {
	return color.Gray16{Y: uint16(red(col))}, nil
}

// Green maps to just the blue component in 8 bits
func Green(col color.Color) (color.Color, error) {
	return color.Gray{Y: uint8(green(col) >> 8)}, nil
}

// Green16 maps to just the blue component in 16 bits
func Green16(col color.Color) (color.Color, error) {
	return color.Gray16{Y: uint16(blue(col))}, nil
}

// Blue maps to just the blue component in 8 bits
func Blue(col color.Color) (color.Color, error) {
	return color.Gray{Y: uint8(blue(col) >> 8)}, nil
}

// Blue16 maps to just the blue component in 16 bits
func Blue16(col color.Color) (color.Color, error) {
	return color.Gray16{Y: uint16(blue(col))}, nil
}
