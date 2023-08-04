package io

import (
	"bytes"
	"golang.org/x/image/tiff"
	"image"
	"io"
)

type TIFF struct{}

func (_ TIFF) Decode(r io.Reader) (image.Image, error) {
	return tiff.Decode(r)
}

func (_ TIFF) DecodeConfig(r io.Reader) (image.Config, error) {
	return tiff.DecodeConfig(r)
}

func (_ TIFF) Encode(w io.Writer, img image.Image) error {
	return tiff.Encode(w, img, nil)
}

func (_ TIFF) EncodeBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := tiff.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (_ TIFF) EncodeFFMPEG(img image.Image) ([]string, error) {
	return nil, nil
}
