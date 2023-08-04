package io

import (
	"errors"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type FFMPeg struct{}

// FFMPegSession handles sending frames to FFMPeg
type FFMPegSession struct {
	encoder Encoder               // The image encoder
	cmd     *exec.Cmd             // The ffmpeg command
	r       *io.PipeReader        // stdin to ffmpeg
	w       *io.PipeWriter        // writer to send images to ffmpeg
	pool    png.EncoderBufferPool // pool of buffers to save on memory allocations
}

type bufferPool struct {
	mutex sync.Mutex
	pool  []*png.EncoderBuffer
}

func (b *bufferPool) Get() *png.EncoderBuffer {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.pool) < 1 {
		return nil
	}

	e := b.pool[0]
	b.pool = b.pool[1:]
	return e
}

func (b *bufferPool) Put(e *png.EncoderBuffer) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.pool = append(b.pool, e)
}

func (_ FFMPeg) new(fileName string, frameRate int, srcArgs []string, encoder Encoder) (*FFMPegSession, error) {
	session := &FFMPegSession{
		pool:    &bufferPool{},
		encoder: encoder,
	}

	frameRateS := strconv.Itoa(frameRate)

	var args []string

	if len(srcArgs) > 0 {
		args = append(args, srcArgs...)
	}

	args = append(args,
		"-y",
		"-framerate", frameRateS,
		"-i", "-", // pipe from stdin
	)

	args = append(args,
		"-r", frameRateS,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
	)

	args = append(args, fileName)

	session.cmd = exec.Command("ffmpeg", args...)

	session.r, session.w = io.Pipe()

	session.cmd.Stdin = session.r

	if log.IsVerbose() {
		session.cmd.Stdout, session.cmd.Stderr = os.Stdout, os.Stderr
	}

	return session, session.cmd.Start()
}

func (f FFMPeg) NewJpg(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.new(fileName, frameRate, nil, &JPEG{})
}

// NewPng creates a new FFMPegSession or an error if ffmpeg could not be started.
func (f FFMPeg) NewPng(fileName string, frameRate int) (*FFMPegSession, error) {
	return f.new(fileName, frameRate, nil, &PNG{})
}

func (f FFMPeg) NewRaw(fileName string, frameRate int, src image.Image) (*FFMPegSession, error) {
	var args []string

	if src == nil {
		return nil, errors.New("image required for raw")
	}

	b := src.Bounds()
	args = append(args,
		"-f", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", b.Dx(), b.Dy()),
	)

	if _, ok := src.(*image.RGBA); ok {
		args = append(args, "-pix_fmt", "rgba")
	} else if _, ok := src.(*image.RGBA64); ok {
		args = append(args, "-pix_fmt", "rgba64")
	} else {
		return nil, fmt.Errorf("unsupported raw image format %T", src)
	}

	return f.new(fileName, frameRate, args, &Raw{})
}

// Close closes the FFMPegSession.
// Normally this is handled by try-resources
func (s *FFMPegSession) Close() error {
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
	return s.w.Write(b)
}

// WriteImage writes an image to ffmpeg.
func (s *FFMPegSession) WriteImage(img image.Image) error {
	b, err := s.encoder.EncodeBytes(img)

	if err == nil {
		_, err = s.Write(b)
	}

	return err
}

func (s *FFMPegSession) EncodeBytes(img image.Image) ([]byte, error) {
	return s.encoder.EncodeBytes(img)
}
