package image

import (
	"fmt"
	"github.com/peter-mount/go-anim/script/render"
	"image"
	"os"
	"strings"
)

// Encoder returns an Encoder appropriate for the specified fileName
func (g *Image) Encoder(fileName string) (render.Encoder, error) {
	for k, h := range g.encoders {
		if strings.HasSuffix(fileName, k) {
			return h, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type %q", fileName)
}

// Decoder returns a Decoder appropriate for the specified fileName
func (g *Image) Decoder(fileName string) (render.Decoder, error) {
	for k, h := range g.decoders {
		if strings.HasSuffix(fileName, k) {
			return h, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type %q", fileName)
}

// WriteImage writes a single image to the specified fileName
func (g *Image) WriteImage(fileName string, img image.Image) error {
	encoder, err := g.Encoder(fileName)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return encoder.Encode(f, img)
}

// ReadImage reads a single image from a file
func (g *Image) ReadImage(fileName string) (image.Image, error) {
	decoder, err := g.Decoder(fileName)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return decoder.Decode(f)
}
