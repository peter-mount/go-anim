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

func (e *encoder) writeStart(w io.Writer) (uint64, error) {
	var offset uint64
	var b bytes.Buffer
	err := e.writeHeader(&b)
	offset = uint64(b.Len())

	if err == nil {
		_, err = b.WriteTo(w)
	}
	return offset, err
}

func (e *encoder) writeHeader(w io.Writer) error {
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

	offset, err := e.writeStart(w)
	if err != nil {
		return err
	}

	return e.encodeImageCompressed(w, m, offset)
}

/*func (e *encoder) encodeImageUncompressed(w io.Writer, m image.Image, offset uint64) error {
	// Point offset to after the chunk headers
	offset += uint64(8 * e.height)

	for y := e.bounds.Min.Y; y <= e.bounds.Max.Y; y++ {
		if err := exr.Write(w, &offset); err != nil {
			return err
		}
		offset += uint64(e.dataSize + 8)
	}

	for y := e.bounds.Min.Y; y <= e.bounds.Max.Y; y++ {
		err := e.writeScanBlock(y, w)
		if err == nil {
			err = e.writeScanline(y, w, m)
		}
		if err != nil {
			return err
		}
	}

	return nil
}*/

func (e *encoder) encodeImageCompressed(w io.Writer, m image.Image, offset uint64) error {
	var compressor exr.Compressor
	switch e.compression {
	case exr.CompressionZIP:
		compressor = exr.NewZipCompressor()
	case exr.CompressionNone:
		compressor = exr.NewNopCompressor()
	default:
		return fmt.Errorf("compression %d unsupported", e.compression)
	}

	lc := e.compression.LineCount()
	fmt.Printf("lineCount %d\n", lc)

	// Compression works on multiple scanlines; so we have to store it in memory first
	var offsets []uint64
	var chunks [][]byte

	chunkSize := uint64(0)

	buffer := &bytes.Buffer{}
	for y := e.bounds.Min.Y; y <= e.bounds.Max.Y; y += lc {
		buffer.Reset()

		for y1 := 0; y1 < lc && (y+y1) <= e.bounds.Max.Y; y1++ {
			if err := e.writeScanline(y+y1, buffer, m); err != nil {
				return err
			}
		}

		cb, err := compressor.Compress(buffer.Bytes())
		if err != nil {
			return err
		}

		// Now set the header
		b := binary.LittleEndian.AppendUint32(nil, uint32(y))
		b = binary.LittleEndian.AppendUint32(b, uint32(len(cb)))
		b = append(b, cb...)

		//fmt.Printf("chunk %d offset %d\n", len(chunks), offset)
		chunks = append(chunks, b)
		offsets = append(offsets, chunkSize)
		chunkSize += uint64(len(b))
	}

	// Now move offsets to include the header & offset table
	offset += uint64(8 * len(offsets))
	for i, o := range offsets {
		offsets[i] = o + offset
	}

	chunkCount := exr.ChunkCount(e.dataWindow, e.compression)
	fmt.Printf("offsets %d chunks %d expected %d\n", len(offsets), len(chunks), chunkCount)

	// Finally write the offsets then the chunks
	if err := exr.Write(w, offsets); err != nil {
		return err
	}

	for _, chunk := range chunks {
		if err := exr.Write(w, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) writeScanBlock(y, lc int, w io.Writer) error {
	// Each line starts with the Y then the size in bytes
	y0 := uint32(y)
	ps := uint32(e.pixelSize)
	err := exr.Write(w, &y0)
	if err == nil {
		err = exr.Write(w, &ps)
	}
	return err
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

// encodeNative encodes an RGBAImage directly
func (e *encoder) encodeNative(w io.Writer, i *RGBAImage) error {
	return fmt.Errorf("unsupported image %T", i)
}

type lineBuffer struct {
}

func (e *encoder) newLineBuffer() {

}
