package io

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type JPEG struct{}

func (_ JPEG) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}

func (_ JPEG) DecodeConfig(r io.Reader) (image.Config, error) {
	return jpeg.DecodeConfig(r)
}

func (_ JPEG) Encode(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
}

func (_ JPEG) EncodeBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
