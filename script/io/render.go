package io

import (
	"fmt"
	"github.com/peter-mount/go-anim/script/util"
	"image"
	"path/filepath"
)

type Render struct{}

func (r Render) New(fileName string, frameRate int) (RenderStream, error) {
	switch filepath.Ext(fileName) {
	// mp4 video
	case ".mp4":
		return r.ffmpeg(fileName, frameRate, &Raw{}), nil

	case ".png":
		return r.frames(fileName, frameRate, &PNG{}), nil

	case ".jpg", ".jpeg":
		return r.frames(fileName, frameRate, &JPEG{}), nil

	case ".tif", ".tiff":
		return r.frames(fileName, frameRate, &TIFF{}), nil

	default:
		return nil, fmt.Errorf("unsupported file type %q", fileName)
	}
}

type RenderStream interface {
	Writer
	TimeCode() *util.TimeCode
	EncodeBytes(img image.Image) ([]byte, error)
}

type RenderStreamBase struct {
	Writer
	fileName string                      // Output fileName
	timeCode *util.TimeCode              // TimeCode
	encoder  Encoder                     // Frame encoder
	init     func(img image.Image) error // init function
	write    func(b []byte) (int, error) // write function
}

func (s *RenderStreamBase) Init(_ image.Image) error {
	return nil
}

func (s *RenderStreamBase) WriteBytes(b []byte) (int, error) {
	if s.init != nil {
		if err := s.init(nil); err != nil {
			return 0, err
		}
	}

	n, err := s.write(b)
	if err == nil {
		s.TimeCode().Next()
	}
	return n, err
}

// WriteImage writes an image to the stream.
func (s *RenderStreamBase) WriteImage(img image.Image) error {
	return s.WriteImageN(img, 1)
}

// WriteImageN writes an image to stream multiple times
func (s *RenderStreamBase) WriteImageN(img image.Image, num int) error {
	if num < 1 {
		return fmt.Errorf("cannot write %d images, must be >=1", num)
	}

	if s.init != nil {
		if err := s.init(img); err != nil {
			return err
		}
	}

	// Encode the frame
	b, err := s.encoder.EncodeBytes(img)

	// Write num copies of the frame
	for n := 0; n < num && err == nil; n++ {
		// Call our write so we increment the TimeCode
		_, err = s.WriteBytes(b)
	}

	return err
}

func (s *RenderStreamBase) EncodeBytes(img image.Image) ([]byte, error) {
	return s.encoder.EncodeBytes(img)
}

func (s *RenderStreamBase) TimeCode() *util.TimeCode {
	return s.timeCode
}

func (s *RenderStreamBase) FrameRate() int {
	return s.TimeCode().FrameRate()
}

func (s *RenderStreamBase) FrameRateF() float64 {
	return s.TimeCode().FrameRateF()
}
