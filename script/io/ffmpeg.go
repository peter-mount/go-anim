package io

import (
	"bytes"
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
	cmd  *exec.Cmd
	r    *io.PipeReader
	w    *io.PipeWriter
	err  error
	pool png.EncoderBufferPool
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

// New creates a new FFMPegSession or an error if ffmpeg could not be started.
func (_ FFMPeg) New(fileName string, frameRate int) (*FFMPegSession, error) {
	session := &FFMPegSession{pool: &bufferPool{}}

	frameRateS := strconv.Itoa(frameRate)

	session.cmd = exec.Command(
		"ffmpeg",
		"-y",
		"-framerate", frameRateS,
		"-i", "-", // pipe from stdin
		"-r", frameRateS,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		fileName,
	)

	session.r, session.w = io.Pipe()

	session.cmd.Stdin = session.r

	if log.IsVerbose() {
		session.cmd.Stdout, session.cmd.Stderr = os.Stdout, os.Stderr
	}

	return session, session.cmd.Start()
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
	var buf bytes.Buffer

	// NoCompression as we are piping to ffmpeg so why spend time compressing?
	// BufferPool, so we reuse buffers to save memory allocations
	enc := png.Encoder{
		CompressionLevel: png.NoCompression,
		BufferPool:       s.pool,
	}

	err := enc.Encode(&buf, img)
	if err == nil {
		_, err = s.Write(buf.Bytes())
	}
	return err
}
