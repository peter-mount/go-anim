package exr

import (
	"bytes"
	"compress/zlib"
)

type Compressor interface {
	Compress(src []byte) ([]byte, error)
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

	// Use Zip level 4 not default 6 to improve performance:
	// https://aras-p.info/blog/2021/08/05/EXR-Zip-compression-levels/
	// https://github.com/AcademySoftwareFoundation/openexr/pull/1125
	w, err := zlib.NewWriterLevelDict(out, 4, nil)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(src); err != nil {
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
