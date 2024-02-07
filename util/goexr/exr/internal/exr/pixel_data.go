package exr

import (
	"fmt"
	"io"

	"github.com/x448/float16"
)

type PixelData interface {
	LineSize() int32
	ReadLine(in io.Reader, y int32) error
	WriteLine(w io.Writer, y int32) error
	Float32(x, y int) float32
	Set(x, y int, v float32)
	PixelType() PixelType
}

func NewNopPixelData(value float32) PixelData {
	return &nopPixelData{
		value: value,
	}
}

type nopPixelData struct {
	value float32
}

func (d *nopPixelData) PixelType() PixelType { return -1 }

func (d *nopPixelData) LineSize() int32 {
	return 0
}

func (d *nopPixelData) ReadLine(in io.Reader, y int32) error {
	return fmt.Errorf("cannot read into nop pixel data")
}

func (d *nopPixelData) WriteLine(w io.Writer, y int32) error {
	return fmt.Errorf("cannot write from nop pixel data")
}

func (d *nopPixelData) Float32(x, y int) float32 {
	return d.value
}

func (d *nopPixelData) Set(x, y int, v float32) {}

func NewUint32PixelData(window Box2i, xSampling, ySampling int32) PixelData {
	return &uint32PixelData{
		width: window.Width() / xSampling,
	}
}

type uint32PixelData struct {
	width int32
}

func (d *uint32PixelData) PixelType() PixelType { return PixelTypeUint }

func (d *uint32PixelData) LineSize() int32 {
	return d.width * 4
}

func (d *uint32PixelData) ReadLine(in io.Reader, y int32) error {
	if _, err := io.CopyN(io.Discard, in, int64(d.width)*4); err != nil {
		return fmt.Errorf("error reading uint32 pixel slice: %w", err)
	}
	return nil
}

func (d *uint32PixelData) WriteLine(w io.Writer, y int32) error {
	return fmt.Errorf("cannot write from nop pixel data")
}

func (d *uint32PixelData) Float32(x, y int) float32 {
	return 0.0 // uint32 is used for object reference, not colors
}

func (d *uint32PixelData) Set(x, y int, v float32) {}

func NewFloat16PixelData(window Box2i, xSampling, ySampling int32) PixelData {
	width := window.Width() / xSampling
	height := window.Height() / ySampling
	return &float16PixelData{
		window:    window,
		xSampling: xSampling,
		ySampling: ySampling,
		pixels:    make([]float16.Float16, width*height),
	}
}

type float16PixelData struct {
	window    Box2i
	xSampling int32
	ySampling int32
	pixels    []float16.Float16
}

func (d *float16PixelData) PixelType() PixelType { return PixelTypeHalf }

func (d *float16PixelData) LineSize() int32 {
	width := d.window.Width() / d.xSampling
	return width * 2
}

func (d *float16PixelData) ReadLine(in io.Reader, y int32) error {
	width := d.window.Width() / d.xSampling
	y = (y - d.window.YMin) / d.ySampling
	offset := y * width
	if err := Read(in, d.pixels[offset:offset+width:offset+width]); err != nil {
		return fmt.Errorf("error reading float16 pixel slice: %w", err)
	}
	return nil
}

func (d *float16PixelData) WriteLine(w io.Writer, y int32) error {
	width := d.window.Width() / d.xSampling
	y = (y - d.window.YMin) / d.ySampling
	offset := y * width
	if err := Write(w, d.pixels[offset:offset+width:offset+width]); err != nil {
		return fmt.Errorf("error writing float16 pixel slice: %w", err)
	}
	return nil
}

func (d *float16PixelData) Float32(x, y int) float32 {
	offX := (int32(x) - d.window.XMin) / d.xSampling
	offY := (int32(y) - d.window.YMin) / d.ySampling
	width := d.window.Width() / d.xSampling

	value := d.pixels[offX+width*offY]
	if value.IsInf(0) {
		value = float16.Frombits(uint16(0x7bff)) // max value
	}
	if value.IsNaN() {
		value = float16.Frombits(uint16(0x0000)) // min value
	}
	return value.Float32()
}

func (d *float16PixelData) Set(x, y int, v float32) {
	offX := (int32(x) - d.window.XMin) / d.xSampling
	offY := (int32(y) - d.window.YMin) / d.ySampling
	width := d.window.Width() / d.xSampling

	value := float16.Fromfloat32(v)
	if value.IsInf(0) {
		value = float16.Frombits(uint16(0x7bff)) // max value
	}
	if value.IsNaN() {
		value = float16.Frombits(uint16(0x0000)) // min value
	}
	d.pixels[offX+width*offY] = value
}

func NewFloat32PixelData(window Box2i, xSampling, ySampling int32) PixelData {
	width := window.Width() / xSampling
	height := window.Height() / ySampling
	return &float32PixelData{
		window:    window,
		xSampling: xSampling,
		ySampling: ySampling,
		pixels:    make([]float32, width*height),
	}
}

type float32PixelData struct {
	window    Box2i
	xSampling int32
	ySampling int32
	pixels    []float32
}

func (d *float32PixelData) PixelType() PixelType { return PixelTypeFloat }

func (d *float32PixelData) LineSize() int32 {
	width := d.window.Width() / d.xSampling
	return width * 4
}

func (d *float32PixelData) ReadLine(in io.Reader, y int32) error {
	width := d.window.Width() / d.xSampling
	y = (y - d.window.YMin) / d.ySampling
	offset := y * width
	if err := Read(in, d.pixels[offset:offset+width:offset+width]); err != nil {
		return fmt.Errorf("error reading float32 pixel slice: %w", err)
	}
	return nil
}

func (d *float32PixelData) WriteLine(w io.Writer, y int32) error {
	width := d.window.Width() / d.xSampling
	y = (y - d.window.YMin) / d.ySampling
	offset := y * width
	if err := Write(w, d.pixels[offset:offset+width:offset+width]); err != nil {
		return fmt.Errorf("error writing float32 pixel slice: %w", err)
	}
	return nil
}

func (d *float32PixelData) Float32(x, y int) float32 {
	offX := (int32(x) - d.window.XMin) / d.xSampling
	offY := (int32(y) - d.window.YMin) / d.ySampling
	width := d.window.Width() / d.xSampling
	return d.pixels[offX+width*offY]
}

func (d *float32PixelData) Set(x, y int, v float32) {
	offX := (int32(x) - d.window.XMin) / d.xSampling
	offY := (int32(y) - d.window.YMin) / d.ySampling
	width := d.window.Width() / d.xSampling

	d.pixels[offX+width*offY] = v
}
