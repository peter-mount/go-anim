package io

import (
	"fmt"
	"image"
	"io"
)

type Raw struct{}

func (r Raw) Encode(w io.Writer, img image.Image) error {
	b, err := r.EncodeBytes(img)
	if err != nil {
		return err
	}

	n, err := w.Write(b)
	if err == nil && n != len(b) {
		return fmt.Errorf("Wrote %d/%d bytes", n, len(b))
	}
	return err
}

func (r Raw) EncodeBytes(img image.Image) ([]byte, error) {
	if src, ok := img.(*image.RGBA); ok {
		return src.Pix, nil
	}
	return nil, fmt.Errorf("unsupported image %T", img)
}
