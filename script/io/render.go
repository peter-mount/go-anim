package io

import (
	"fmt"
	"github.com/peter-mount/go-anim/script/util"
	"github.com/peter-mount/go-script/packages"
	"image"
	"strings"
)

func init() {
	r := &Render{}

	// Populate the extensions.
	// This is first come, first served so ensure that the longer
	// variants are first, e.g. .raw.mp4 is before .mp4
	r.renderers = []rendererHandler{
		// .mp4 frame types
		{suffix: ".raw.mp4", handler: r.newRawMp4},
		{suffix: ".png.mp4", handler: r.newPngMp4},
		{suffix: ".jpg.mp4", handler: r.newJpegMp4},
		{suffix: ".jpeg.mp4", handler: r.newJpegMp4},
		{suffix: ".tiff.mp4", handler: r.newTiffMp4},
		{suffix: ".tif.mp4", handler: r.newTiffMp4},
		// .mp4 default using raw frames
		{suffix: ".mp4", handler: r.newRawMp4},
		// directory frame types
		{suffix: ".png", handler: r.newPng},
		{suffix: ".jpg", handler: r.newJpeg},
		{suffix: ".jpeg", handler: r.newJpeg},
		{suffix: ".tiff", handler: r.newTiff},
		{suffix: ".tif", handler: r.newTiff},
		// tar frame types
		{suffix: ".png.tar", handler: r.newPngTar},
		{suffix: ".jpg.tar", handler: r.newJpegTar},
		{suffix: ".jpeg.tar", handler: r.newJpegTar},
		{suffix: ".tiff.tar", handler: r.newTiffTar},
		{suffix: ".tif.tar", handler: r.newTiffTar},
		// tar default using png frames
		{suffix: ".tar", handler: r.newPngTar},
	}

	packages.Register("render", r)
}

type Render struct {
	renderers []rendererHandler
}
type rendererHandler struct {
	suffix  string
	handler func(fileName string, frameRate int) RenderStream
}

func (r Render) New(fileName string, frameRate int) (RenderStream, error) {
	for _, h := range r.renderers {
		if strings.HasSuffix(fileName, h.suffix) {
			return h.handler(fileName, frameRate), nil
		}
	}

	return nil, fmt.Errorf("unsupported file type %q", fileName)
}

func (r Render) newRawMp4(fileName string, frameRate int) RenderStream {
	return r.ffmpeg(fileName, frameRate, &Raw{})
}

func (r Render) newPngMp4(fileName string, frameRate int) RenderStream {
	return r.ffmpeg(fileName, frameRate, &PNG{})
}

func (r Render) newJpegMp4(fileName string, frameRate int) RenderStream {
	return r.ffmpeg(fileName, frameRate, &JPEG{})
}

func (r Render) newTiffMp4(fileName string, frameRate int) RenderStream {
	return r.ffmpeg(fileName, frameRate, &TIFF{})
}

func (r Render) newPng(fileName string, frameRate int) RenderStream {
	return r.frames(fileName, frameRate, &PNG{})
}

func (r Render) newJpeg(fileName string, frameRate int) RenderStream {
	return r.frames(fileName, frameRate, &JPEG{})
}

func (r Render) newTiff(fileName string, frameRate int) RenderStream {
	return r.frames(fileName, frameRate, &TIFF{})
}

func (r Render) newPngTar(fileName string, frameRate int) RenderStream {
	return r.tar(fileName, frameRate, &PNG{}, ".png")
}

func (r Render) newJpegTar(fileName string, frameRate int) RenderStream {
	return r.tar(fileName, frameRate, &JPEG{}, ".jpg")
}

func (r Render) newTiffTar(fileName string, frameRate int) RenderStream {
	return r.tar(fileName, frameRate, &TIFF{}, ".tiff")
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
