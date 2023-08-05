package io

import (
	"fmt"
	"github.com/peter-mount/go-anim/script/util"
	"github.com/peter-mount/go-kernel/v2/log"
	"image"
	"io"
	"os"
	"os/exec"
	"strconv"
)

type FFMPeg struct{}

// FFMPegSession handles sending frames to FFMPeg
type FFMPegSession struct {
	fileName string              // Output fileName
	timeCode *util.TimeCode      // TimeCode
	encoder  FFMPegSessionSource // The image encoder
	cmd      *exec.Cmd           // The ffmpeg command
	r        *io.PipeReader      // stdin to ffmpeg
	w        *io.PipeWriter      // writer to send images to ffmpeg
}

type FFMPegSessionSource interface {
	Encoder
	EncodeFFMPEG(img image.Image) ([]string, error)
}

func (_ FFMPeg) new(fileName string, frameRate int, encoder FFMPegSessionSource) *FFMPegSession {
	return &FFMPegSession{
		fileName: fileName,
		encoder:  encoder,
		timeCode: util.NewTimeCode(frameRate),
	}
}

func (s *FFMPegSession) init(img image.Image) error {

	// Get initial source parameters from the encoder, may be nil if none required
	args, err := s.encoder.EncodeFFMPEG(img)
	if err != nil {
		return err
	}

	frameRateS := strconv.Itoa(s.TimeCode().FrameRate())

	args = append(args,
		// Required source parameters
		"-y",
		"-framerate", frameRateS,
		// pipe from stdin
		"-i", "-",
		// Always provide the start time code
		"-timecode", s.TimeCode().StartTimeCode().TimeCode(),
		// Now the destination parameters
		"-r", frameRateS,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		// Destination file name
		s.fileName,
	)

	s.cmd = exec.Command("ffmpeg", args...)

	s.r, s.w = io.Pipe()

	s.cmd.Stdin = s.r

	if log.IsVerbose() {
		s.cmd.Stdout, s.cmd.Stderr = os.Stdout, os.Stderr
	}

	return s.cmd.Start()
}

// New creates a new FFMPegSession. The default is to use raw frame formatting
func (f FFMPeg) New(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.NewRaw(fileName, frameRate)
}

// NewJpg creates a new FFMPegSession using JPG as the frame format
func (f FFMPeg) NewJpg(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.new(fileName, frameRate, &JPEG{}), nil
}

// NewPng creates a new FFMPegSession using PNG as the frame format
func (f FFMPeg) NewPng(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.new(fileName, frameRate, &PNG{}), nil
}

// NewRaw creates a new FFMPegSession using the raw image format of the first frame
func (f FFMPeg) NewRaw(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.new(fileName, frameRate, &Raw{}), nil
}

// Close closes the FFMPegSession.
// Normally this is handled by try-resources
func (s *FFMPegSession) Close() error {
	if s.w == nil {
		return nil
	}

	if err := s.w.Close(); err != nil {
		fmt.Println("Error closing ffmpeg stream", err)
		_ = s.cmd.Process.Kill()
	}

	return s.cmd.Wait()
}

// Write a block of bytes to ffmpeg.
// Normally this will be a pre-encoded image - e.g. when a frame is used
// multiple times, only render it once
func (s *FFMPegSession) Write(b []byte) (int, error) {
	// Lazy init ffmpeg, passing nil. If the encoder requires an image it should fail
	if s.cmd == nil {
		if err := s.init(nil); err != nil {
			return 0, err
		}
	}

	n, err := s.w.Write(b)
	if err == nil {
		s.TimeCode().Next()
	}
	return n, err
}

// WriteImage writes an image to ffmpeg.
func (s *FFMPegSession) WriteImage(img image.Image) error {
	return s.WriteImageMulti(img, 1)
}

// WriteImageMulti writes an image to ffmpeg multiple times
func (s *FFMPegSession) WriteImageMulti(img image.Image, num int) error {
	if num < 1 {
		return fmt.Errorf("cannot write %d images, must be >=1", num)
	}

	// Lazy init ffmpeg passing the image
	if s.cmd == nil {
		if err := s.init(img); err != nil {
			return err
		}
	}

	// Encode the frame
	b, err := s.encoder.EncodeBytes(img)

	// Write num copies of the frame
	for n := 0; n < num && err == nil; n++ {
		// Call our write so we increment the TimeCode
		_, err = s.Write(b)
	}

	return err
}

func (s *FFMPegSession) EncodeBytes(img image.Image) ([]byte, error) {
	return s.encoder.EncodeBytes(img)
}

func (s *FFMPegSession) TimeCode() *util.TimeCode {
	return s.timeCode
}
