package exr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/peter-mount/go-anim/util/goexr/exr/attributes"
	"github.com/peter-mount/go-anim/util/goexr/exr/internal/exr"
	"github.com/x448/float16"
	"image"
	"io"
	"math"
)

func Encode(w io.Writer, m image.Image) error {
	return NewEncoder().Encode(w, m)
}

type Encoder interface {
	Encode(w io.Writer, m image.Image) error
	Compress(b bool) Encoder
	Float16() Encoder
	Float32() Encoder
}

type encoder struct {
	XSampling     int32
	YSampling     int32
	dataWindow    exr.Box2i
	displayWindow exr.Box2i
	lineOrder     exr.LineOrder
	pixelType     exr.PixelType
	channels      exr.ChannelList
	compression   exr.Compression

	bounds                 image.Rectangle
	width, height          int
	pixelWidth             int
	lineSize               int
	pixelSize              int
	dataSize               int32
	lineBuffer             []byte
	aOff, bOff, gOff, rOff int
}

func NewEncoder() Encoder {
	return &encoder{
		pixelType:   exr.PixelTypeHalf,
		compression: exr.CompressionNone,
	}
}

func (e *encoder) Float16() Encoder {
	e.pixelType = exr.PixelTypeHalf
	return e
}

func (e *encoder) Float32() Encoder {
	e.pixelType = exr.PixelTypeFloat
	return e
}

func (e *encoder) Compress(b bool) Encoder {
	if b {
		e.compression = exr.CompressionZIP
	} else {
		e.compression = exr.CompressionNone
	}
	return e
}

func (e *encoder) Encode(w io.Writer, m image.Image) error {
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

	if ni, ok := m.(*RGBAImage); ok {
		return e.encodeNative(w, ni)
	}

	return e.encodeImage(w, m)
}

func (e *encoder) writeStart(w io.Writer, m image.Image) (uint64, error) {
	var offset uint64
	var b bytes.Buffer
	err := e.writeHeader(&b, m)
	offset = uint64(b.Len())

	if err == nil {
		_, err = b.WriteTo(w)
	}
	return offset, err
}

func (e *encoder) writeHeader(w io.Writer, m image.Image) error {
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
		err = exr.WriteAttribute(w, exr.AttributeNameCompression, exr.AttributeTypeCompression, e.compression)
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

	// Include any additional attributes
	if attrs, ok := m.(attributes.ImageAttributes); ok {
		err = attrs.ForEach(func(a attributes.Attribute) error {
			return exr.WriteAttributeBytes(w, exr.AttributeName(a.Name), exr.AttributeType(a.Type), int32(len(a.Data)), a.Data)
		})
	}

	// Terminate the header
	if err == nil {
		_, err = w.Write([]byte{0x00})
	}

	return err
}

var (
	// The Channels to render, must be in Alphabetical order as that's the order defined for a scan line within exr.
	components = []string{"A", "B", "G", "R"}
)

// encodeImage encodes a normal Image. RGBAImage is handled directly by encodeNative
func (e *encoder) encodeImage(w io.Writer, m image.Image) error {

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

	// As these images are usually defined with a pixel array rather than by each component
	// we need to handle this by component and then line by line
	e.bounds = m.Bounds()
	e.width = e.bounds.Max.X - e.bounds.Min.X + 1
	e.height = e.bounds.Max.Y - e.bounds.Min.Y + 1

	switch e.pixelType {
	case exr.PixelTypeHalf:
		e.pixelWidth = 2
	case exr.PixelTypeFloat, exr.PixelTypeUint:
		e.pixelWidth = 4
	default:
		return fmt.Errorf("unknown PixelType %d", e.pixelType)
	}

	// Each row is y, line size * channelCount followed by that number of bytes
	e.lineSize = e.width * e.pixelWidth
	e.pixelSize = len(e.channels) * e.lineSize
	e.dataSize = int32(e.pixelSize)
	e.lineBuffer = make([]byte, e.dataSize)

	e.aOff = 0
	e.bOff = e.aOff + e.lineSize
	e.gOff = e.bOff + e.lineSize
	e.rOff = e.gOff + e.lineSize

	offset, err := e.writeStart(w, m)
	if err != nil {
		return err
	}

	cw := exr.NewChunkWriter(e.compression, w, offset, e.lineSize, len(e.channels), e.bounds, e.dataWindow)

	return cw.Write(func(y int, w io.Writer) error {
		return e.writeScanline(y, w, m)
	})
}

func (e *encoder) writeScanline(y int, w io.Writer, m image.Image) error {
	off := 0
	for x := e.bounds.Min.X; x <= e.bounds.Max.X; x++ {
		r, g, b, a := m.At(x, y).RGBA()
		switch e.pixelType {
		case exr.PixelTypeHalf:
			binary.LittleEndian.PutUint16(e.lineBuffer[e.aOff+off:e.aOff+off+2], float16.Fromfloat32(float32(a)/float32(0xffff)).Bits())
			binary.LittleEndian.PutUint16(e.lineBuffer[e.bOff+off:e.bOff+off+2], float16.Fromfloat32(float32(b)/float32(0xffff)).Bits())
			binary.LittleEndian.PutUint16(e.lineBuffer[e.gOff+off:e.gOff+off+2], float16.Fromfloat32(float32(g)/float32(0xffff)).Bits())
			binary.LittleEndian.PutUint16(e.lineBuffer[e.rOff+off:e.rOff+off+2], float16.Fromfloat32(float32(r)/float32(0xffff)).Bits())
		case exr.PixelTypeFloat:
			binary.LittleEndian.PutUint32(e.lineBuffer[e.aOff+off:e.aOff+off+4], math.Float32bits(float32(a&0xffff)/float32(0xffff)))
			binary.LittleEndian.PutUint32(e.lineBuffer[e.bOff+off:e.bOff+off+4], math.Float32bits(float32(b&0xffff)/float32(0xffff)))
			binary.LittleEndian.PutUint32(e.lineBuffer[e.gOff+off:e.gOff+off+4], math.Float32bits(float32(g&0xffff)/float32(0xffff)))
			binary.LittleEndian.PutUint32(e.lineBuffer[e.rOff+off:e.rOff+off+4], math.Float32bits(float32(r&0xffff)/float32(0xffff)))
		case exr.PixelTypeUint:
			binary.LittleEndian.PutUint32(e.lineBuffer[e.aOff+off:e.aOff+off+4], a<<16)
			binary.LittleEndian.PutUint32(e.lineBuffer[e.bOff+off:e.bOff+off+4], b<<16)
			binary.LittleEndian.PutUint32(e.lineBuffer[e.gOff+off:e.gOff+off+4], g<<16)
			binary.LittleEndian.PutUint32(e.lineBuffer[e.rOff+off:e.rOff+off+4], r<<16)
		}
		off += e.pixelWidth
	}

	_, err := w.Write(e.lineBuffer)
	return err
}

// encodeNative encodes an RGBAImage directly as we keep the channels in the correct format already so
// this will be more performant than encodeImage
func (e *encoder) encodeNative(w io.Writer, i *RGBAImage) error {

	// Add the fixed RGBA channels - note: although here the order doesn't matter,
	// keep them alphabetical as we use this order for ordering the channels for each row
	// and there they must be in this order!
	e.channels = nil
	e.addChannel("A", i.channelA)
	e.addChannel("B", i.channelB)
	e.addChannel("G", i.channelG)
	e.addChannel("R", i.channelR)

	offset, err := e.writeStart(w, i)
	if err != nil {
		return err
	}

	cw := exr.NewChunkWriter(e.compression, w, offset, e.lineSize, len(e.channels), i.Bounds(), e.dataWindow)

	return cw.Write(func(y int, w io.Writer) error {
		var err error
		for _, ch := range e.channels {
			switch ch.Name {
			case "A":
				err = i.channelA.WriteLine(w, int32(y))
			case "B":
				err = i.channelB.WriteLine(w, int32(y))
			case "G":
				err = i.channelG.WriteLine(w, int32(y))
			case "R":
				err = i.channelR.WriteLine(w, int32(y))
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (e *encoder) addChannel(n string, pd exr.PixelData) {
	switch pd.PixelType() {
	case exr.PixelTypeHalf, exr.PixelTypeFloat:
		e.channels = append(e.channels, exr.Channel{
			Name:      n,
			PixelType: e.pixelType,
			Linear:    true,
			XSampling: e.XSampling,
			YSampling: e.YSampling,
		})
	}
}
