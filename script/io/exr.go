package io

import (
	"bytes"
	"github.com/peter-mount/go-anim/util/goexr/exr"
	"image"
	"io"
)

type EXR struct{}

func (_ EXR) Encoder() RawEncoder {
	return exr.NewEncoder()
}

func (_ EXR) Encode(w io.Writer, img image.Image) error {
	return exr.Encode(w, img)
}

func (_ EXR) EncodeBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := exr.Encode(&buf, img); err != nil {
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
