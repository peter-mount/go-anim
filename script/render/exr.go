package render

import (
	"bytes"
	"errors"
	"github.com/peter-mount/go-anim/util/goexr/exr"
	"image"
	"io"
)

type EXR struct {
	encoder exr.Encoder
}

func (_ EXR) Encoder() RawEncoder {
	return &EXR{encoder: exr.NewEncoder()}
}

func (e EXR) Encode(w io.Writer, img image.Image) error {
	if e.encoder == nil {
		return exr.Encode(w, img)
	} else {
		return e.encoder.Encode(w, img)
	}
}

func (e EXR) EncodeBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (_ EXR) EncodeFFMPEG(_ image.Image) ([]string, error) {
	return nil, nil
}

func (_ EXR) Decode(r io.Reader) (image.Image, error) {
	return exr.Decode(r)
}

func (e EXR) assertEncoder() error {
	if e.encoder == nil {
		return errors.New("unsupported")
	}
	return nil
}

func (e EXR) Compress(b bool) (EXR, error) {
	if err := e.assertEncoder(); err != nil {
		return e, err
	}
	e.encoder.Compress(b)
	return e, nil
}

func (e EXR) Float16() (EXR, error) {
	if err := e.assertEncoder(); err != nil {
		return e, err
	}
	e.encoder.Float16()
	return e, nil
}

func (e EXR) Float32() (EXR, error) {
	if err := e.assertEncoder(); err != nil {
		return e, err
	}
	e.encoder.Float32()
	return e, nil
}
