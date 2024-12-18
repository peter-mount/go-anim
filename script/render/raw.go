package render

import (
	"errors"
	"fmt"
	"image"
	"io"
)

// Raw package handles raw images.
// Currently, RGBA and RGBA64 are supported as those are the most common and the way
// they are stored in memory is directly supported by ffmpeg.
type Raw struct{}

func (r Raw) Decode(_ io.Reader) (image.Image, error) {
	return nil, errors.New("not implemented")
}

func (r Raw) Encoder() RawEncoder {
	return r
}

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

	if src, ok := img.(*image.RGBA64); ok {
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
