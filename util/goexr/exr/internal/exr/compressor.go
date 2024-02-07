package exr

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

type Compressor interface {
	Compress(src []byte) ([]byte, error)
}

func NewCompressor(c Compression) (Compressor, error) {
	switch c {
	case CompressionZIP:
		return NewZipCompressor(), nil
	case CompressionNone:
		return NewNopCompressor(), nil
	default:
		return nil, fmt.Errorf("compression %d unsupported", c)
	}
}

func NewNopCompressor() Compressor {
	return &nopCompressor{}
}

type nopCompressor struct{}

func (d *nopCompressor) Compress(src []byte) ([]byte, error) {
	return src, nil
}

func NewZipCompressor() Compressor {
	return &zipCompressor{}
}

type zipCompressor struct{}

func (d *zipCompressor) Compress(src []byte) ([]byte, error) {
	out := &bytes.Buffer{}

	// interleave scalar
	result := make([]byte, len(src))
	i1 := 0
	i2 := (len(src) + 1) / 2
	j := 0
	for j < len(result) {
		result[i1] = src[j]
		i1++
		j++

		if j < len(result) {
			result[i2] = src[j]
			i2++
			j++
		}
	}

	// delta encode
	p := int(result[0])
	for i := 1; i < len(result); i++ {
		v := int(result[i]) - p + 128 + 256
		p = int(result[i])
		result[i] = byte(v)
	}

	// Use Zip level 4 not default 6 to improve performance:
	// https://aras-p.info/blog/2021/08/05/EXR-Zip-compression-levels/
	// https://github.com/AcademySoftwareFoundation/openexr/pull/1125
	w, err := zlib.NewWriterLevelDict(out, 4, nil)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(result); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	dst := out.Bytes()

	// Note: as per the exr spec, we return the smallest of src or data
	if len(dst) < len(src) {
		return dst, nil
	}
	return src, nil
}

type ChunkWriter struct {
	compression   Compression
	bounds        image.Rectangle
	lc            int
	initialOffset uint64
	offset        uint64
	offsets       []uint64
	chunks        [][]byte
	lineSize      int
	channels      int
	chunkCount    int
	w             io.Writer
	buffer        bytes.Buffer
}

type ChunkWriterHandler func(y int, w io.Writer) error

func NewChunkWriter(c Compression, w io.Writer, offset uint64, lineSize, channels int, b image.Rectangle, dataWindow Box2i) *ChunkWriter {
	chunkCount := ChunkCount(dataWindow, c)
	return &ChunkWriter{
		compression:   c,
		bounds:        b,
		chunkCount:    chunkCount,
		channels:      channels,
		initialOffset: offset,
		offset:        offset + uint64(8*chunkCount),
		lineSize:      lineSize,
		w:             w,
	}
}

func (cw *ChunkWriter) Write(f ChunkWriterHandler) error {
	compressor, err := NewCompressor(cw.compression)
	if err != nil {
		return err
	}

	cw.lc = cw.compression.LineCount()

	if cw.compression == CompressionNone {
		return cw.writeUncompressed(f)
	}
	return cw.writeCompressed(compressor, f)
}

func (cw *ChunkWriter) writeCompressed(compressor Compressor, f ChunkWriterHandler) error {
	for y := cw.bounds.Min.Y; y <= cw.bounds.Max.Y; y += cw.lc {

		cw.buffer.Reset()

		for y1 := 0; y1 < cw.lc && (y+y1) <= cw.bounds.Max.Y; y1++ {
			if err := f(y+y1, &cw.buffer); err != nil {
				return err
			}
		}

		cb, err := compressor.Compress(cw.buffer.Bytes())
		if err != nil {
			return err
		}

		b := cw.prefixHeader(y, cb)
		cw.chunks = append(cw.chunks, b)
		cw.offsets = append(cw.offsets, cw.offset)
		cw.offset += uint64(len(b))
	}

	// Finally write the offsets then the chunks
	if err := Write(cw.w, cw.offsets); err != nil {
		return err
	}

	for _, chunk := range cw.chunks {
		if err := Write(cw.w, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (cw *ChunkWriter) writeUncompressed(f ChunkWriterHandler) error {
	// Create and write the offsets now
	off := cw.offset
	dataSize := uint64((cw.lineSize * cw.channels) + 8)
	for i := 0; i < cw.chunkCount; i++ {
		cw.offsets = append(cw.offsets, off)
		off += dataSize
	}

	err := Write(cw.w, cw.offsets)
	if err != nil {
		return err
	}

	// Now write each line one by one as each one is a chunk in itself
	for y := cw.bounds.Min.Y; y <= cw.bounds.Max.Y; y++ {
		cw.buffer.Reset()

		if err = f(y, &cw.buffer); err != nil {
			return err
		}

		// Now set the header
		b := cw.prefixHeader(y, cw.buffer.Bytes())
		if err = Write(cw.w, b); err != nil {
			return err
		}
	}

	return nil
}

func (cw *ChunkWriter) prefixHeader(y int, cb []byte) []byte {
	b := binary.LittleEndian.AppendUint32(nil, uint32(y))
	b = binary.LittleEndian.AppendUint32(b, uint32(len(cb)))
	return append(b, cb...)
}
