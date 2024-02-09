package io

import (
	"image"
	"io"
	"os"
)

type RawEncoder interface {
	Encode(w io.Writer, img image.Image) error
}

type Encoder interface {
	RawEncoder
	Encoder() RawEncoder
	EncodeBytes(img image.Image) ([]byte, error)
	EncodeFFMPEG(img image.Image) ([]string, error)
}

type Decoder interface {
	Decode(r io.Reader) (image.Image, error)
}

type imageCodec interface {
	Encoder
	Decoder
}

type ImageCodec interface {
	imageCodec
	Read(fileName string) (image.Image, error)
	Write(string, image.Image) error
}

type imageCodecImpl struct {
	codec imageCodec
}

func (c *imageCodecImpl) Encoder() RawEncoder {
	return c.codec.Encoder()
}

func (c *imageCodecImpl) Encode(w io.Writer, img image.Image) error {
	return c.codec.Encode(w, img)
}

func (c *imageCodecImpl) EncodeBytes(img image.Image) ([]byte, error) {
	return c.codec.EncodeBytes(img)
}

func (c *imageCodecImpl) Decode(r io.Reader) (image.Image, error) {
	return c.codec.Decode(r)
}

func (c *imageCodecImpl) Read(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return c.Decode(f)
}

func (c *imageCodecImpl) Write(fileName string, img image.Image) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Encode(f, img)
}

func (c *imageCodecImpl) EncodeFFMPEG(img image.Image) ([]string, error) {
	return c.codec.EncodeFFMPEG(img)
}

func codec(codec imageCodec) ImageCodec {
	return &imageCodecImpl{codec: codec}
}
