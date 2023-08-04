package io

import (
	"errors"
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

func (_ Raw) EncodeFFMPEG(img image.Image) ([]string, error) {
	if img == nil {
		return nil, errors.New("image required")
	}

	var args []string

	b := img.Bounds()
	args = append(args,
		"-f", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", b.Dx(), b.Dy()),
	)

	if _, ok := img.(*image.RGBA); ok {
		args = append(args, "-pix_fmt", "rgba")
	} else if _, ok := img.(*image.RGBA64); ok {
		args = append(args, "-pix_fmt", "rgba64")
	} else {
		return nil, fmt.Errorf("unsupported raw image format %T", img)
	}

	return args, nil
}
