package render

import (
	"fmt"
	"image"
	"os"
	"strings"
)

// Encoder returns an Encoder appropriate for the specified fileName
func (r Render) Encoder(fileName string) (Encoder, error) {
	for k, h := range r.encoders {
		if strings.HasSuffix(fileName, k) {
			return h, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type %q", fileName)
}

// Decoder returns a Decoder appropriate for the specified fileName
func (r Render) Decoder(fileName string) (Decoder, error) {
	for k, h := range r.decoders {
		if strings.HasSuffix(fileName, k) {
			return h, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type %q", fileName)
}

// WriteImage writes a single image to the specified fileName
func (r Render) WriteImage(fileName string, img image.Image) error {
	encoder, err := r.Encoder(fileName)
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
func (r Render) ReadImage(fileName string) (image.Image, error) {
	decoder, err := r.Decoder(fileName)
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
