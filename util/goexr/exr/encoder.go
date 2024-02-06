package exr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/peter-mount/go-anim/util/goexr/exr/internal/exr"
	"github.com/x448/float16"
	"image"
	"io"
	"math"
)

func Encode(w io.Writer, m image.Image) error {
	var e Encoder
	return e.Encode(w, m)
}

type Encoder struct {
	XSampling     int32
	YSampling     int32
	dataWindow    exr.Box2i
	displayWindow exr.Box2i
	lineOrder     exr.LineOrder
	pixelType     exr.PixelType
	channels      exr.ChannelList
}

func (e *Encoder) Encode(w io.Writer, m image.Image) error {
	if e.XSampling < 1 {
		e.XSampling = 1
	}

	if e.YSampling < 1 {
		e.YSampling = 1
	}

	rect := m.Bounds()
	e.dataWindow = exr.Box2iFromRect(rect)
	e.displayWindow = exr.Box2iFromRect(rect)
	e.lineOrder = exr.LineOrderIncreasingY

	e.pixelType = exr.PixelTypeHalf
	//e.pixelType = exr.PixelTypeFloat

	// For now, we have the fixed R, G, B & A channels
	e.channels = nil
	for _, n := range components {
		e.channels = append(e.channels, exr.Channel{
			Name:      n,
			PixelType: e.pixelType,
			Linear:    true,
			XSampling: e.XSampling,
			YSampling: e.YSampling,
		})
	}

	if ni, ok := m.(*RGBAImage); ok {
		return e.encodeNative(w, ni)
	}

	return e.encodeImage(w, m)
}

func (e *Encoder) writeStart(w io.Writer) (uint64, error) {
	var offset uint64
	var b bytes.Buffer
	err := e.writeHeader(&b)
	offset = uint64(b.Len())

	if err == nil {
		_, err = b.WriteTo(w)
	}
	return offset, err
}

func (e *Encoder) writeHeader(w io.Writer) error {
	err := exr.WriteMagic(w)

	if err == nil {
		var version exr.Version
		version = exr.SupportedVersion
		err = exr.Write(w, &version)
	}

	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameChannels, exr.AttributeTypeChannelList, &e.channels)
	}

	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameCompression, exr.AttributeTypeCompression, exr.CompressionNone)
	}

	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameDataWindow, exr.AttributeTypeBox2i, &e.dataWindow)
	}

	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameDisplayWindow, exr.AttributeTypeBox2i, &e.displayWindow)
	}

	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameLineOrder, exr.AttributeTypeLineOrder, &e.lineOrder)
	}

	// Pixel Aspect Ratio - expect 1.0
	if err == nil {
		par := float32(1.0)
		err = exr.WriteAttribute(w, exr.AttributeNamePixelAspectRatio, exr.AttributeTypeFloat, &par)
	}

	// as per https://openexr.com/en/latest/StandardAttributes.html
	// set Width to 1 and center to 0,0 but these are required attributes
	if err == nil {
		err = exr.WriteAttribute(w, exr.AttributeNameScreenWindowCenter, exr.AttributeTypeV2f, &exr.V2F{})
	}
	if err == nil {
		par := float32(1.0)
		err = exr.WriteAttribute(w, exr.AttributeNameScreenWindowWidth, exr.AttributeTypeFloat, &par)
	}

	// Terminate the header
	if err == nil {
		_, err = w.Write([]byte{0x00})
	}

	return err
}

var (
	// The Channels to render, must be in RGBA order.
	// TODO for some reason Alpha breaks the output so we do not include it for now until I understand how this is meant to behave
	components = []string{"R", "G", "B" /*, "A"*/}
)

// encodeImage encodes a normal Image. RGBAImage is handled directly by encodeNative
func (e *Encoder) encodeImage(w io.Writer, m image.Image) error {
	// As these images are usually defined with a pixel array rather than by each component
	// we need to handle this by component and then line by line
	bounds := m.Bounds()
	width := bounds.Max.X - bounds.Min.X + 1
	height := bounds.Max.Y - bounds.Min.Y + 1

	var pixelWidth int
	switch e.pixelType {
	case exr.PixelTypeHalf:
		pixelWidth = 2
	case exr.PixelTypeFloat, exr.PixelTypeUint:
		pixelWidth = 4
	default:
		return fmt.Errorf("unknown PixelType %d", e.pixelType)
	}

	// Each row is y, line size * channelCount followed by that number of bytes
	lineSize := width * pixelWidth
	pixelSize := len(e.channels) * lineSize
	dataSize := int32(pixelSize + 8)
	lineBuffer := make([]byte, dataSize)

	rOff := 8
	gOff := rOff + lineSize
	bOff := gOff + lineSize
	//aOff := bOff + lineSize

	offset, err := e.writeStart(w)
	if err != nil {
		return err
	}

	// Point offset to after the chunk headers
	offset += uint64(8 * height)

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		if err = exr.Write(w, &offset); err != nil {
			return err
		}
		offset += uint64(dataSize)
	}

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		// Each line starts with the Y then the size in bytes
		binary.LittleEndian.PutUint32(lineBuffer[0:4], uint32(y))
		binary.LittleEndian.PutUint32(lineBuffer[4:8], uint32(pixelSize))

		off := 0
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			switch e.pixelType {
			case exr.PixelTypeHalf:
				binary.LittleEndian.PutUint16(lineBuffer[rOff+off:rOff+off+2], float16.Fromfloat32(float32(r)/float32(0xffff)).Bits())
				binary.LittleEndian.PutUint16(lineBuffer[gOff+off:gOff+off+2], float16.Fromfloat32(float32(g)/float32(0xffff)).Bits())
				binary.LittleEndian.PutUint16(lineBuffer[bOff+off:bOff+off+2], float16.Fromfloat32(float32(b)/float32(0xffff)).Bits())
				//binary.LittleEndian.PutUint16(lineBuffer[aOff+off:aOff+off+2], float16.Fromfloat32(float32(a)/float32(0xffff)).Bits())
			case exr.PixelTypeFloat:
				binary.LittleEndian.PutUint32(lineBuffer[rOff+off:rOff+off+4], math.Float32bits(float32(r&0xffff)/float32(0xffff)))
				binary.LittleEndian.PutUint32(lineBuffer[gOff+off:gOff+off+4], math.Float32bits(float32(g&0xffff)/float32(0xffff)))
				binary.LittleEndian.PutUint32(lineBuffer[bOff+off:bOff+off+4], math.Float32bits(float32(b&0xffff)/float32(0xffff)))
				//binary.LittleEndian.PutUint32(lineBuffer[aOff+off:aOff+off+4], math.Float32bits(float32(a&0xffff)/float32(0xffff)))
			case exr.PixelTypeUint:
				binary.LittleEndian.PutUint32(lineBuffer[rOff+off:rOff+off+4], r<<16)
				binary.LittleEndian.PutUint32(lineBuffer[gOff+off:gOff+off+4], g<<16)
				binary.LittleEndian.PutUint32(lineBuffer[bOff+off:bOff+off+4], b<<16)
				//binary.LittleEndian.PutUint32(lineBuffer[aOff+off:aOff+off+4], a<<16)
			}
			off += pixelWidth
		}

		_, err = w.Write(lineBuffer)
		if err != nil {
			return err
		}
	}

	return nil
}

// encodeNative encodes an RGBAImage directly
func (e *Encoder) encodeNative(w io.Writer, i *RGBAImage) error {
	return fmt.Errorf("unsupported image %T", i)
}
