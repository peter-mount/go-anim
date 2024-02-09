package io

import (
	"bytes"
	"image"
	"image/png"
	"io"
)

type PNG struct{}

func (_ PNG) Encoder() RawEncoder {
	return &png.Encoder{}
}

func (_ PNG) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

func (_ PNG) DecodeConfig(r io.Reader) (image.Config, error) {
	return png.DecodeConfig(r)
}

func (_ PNG) Encode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func (_ PNG) EncodeBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (_ PNG) EncodeFFMPEG(img image.Image) ([]string, error) {
	return nil, nil
}
