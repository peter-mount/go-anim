package io

import (
	"image"
	"io"
)

type Encoder interface {
	Encode(w io.Writer, img image.Image) error
	EncodeBytes(img image.Image) ([]byte, error)
}
