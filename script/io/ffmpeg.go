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

// FFMPegSession handles sending frames to FFMPeg
type FFMPegSession struct {
	RenderStreamBase
	encoder FFMPegSessionSource // The image encoder
	cmd     *exec.Cmd           // The ffmpeg command
	r       *io.PipeReader      // stdin to ffmpeg
	w       *io.PipeWriter      // writer to send images to ffmpeg
}

type FFMPegSessionSource interface {
	Encoder
	EncodeFFMPEG(img image.Image) ([]string, error)
}

func (_ Render) ffmpeg(fileName string, frameRate int, encoder FFMPegSessionSource) RenderStream {
	s := &FFMPegSession{
		RenderStreamBase: RenderStreamBase{
			fileName: fileName,
			timeCode: util.NewTimeCode(frameRate),
			encoder:  encoder,
		},
		encoder: encoder,
	}
	s.RenderStreamBase.init = s.init
	s.RenderStreamBase.write = s.writeBytes
	return s
}

func (s *FFMPegSession) init(img image.Image) error {
	if s.cmd != nil {
		return nil
	}

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

	fmt.Println("args", args)
	s.cmd = exec.Command("ffmpeg", args...)

	s.r, s.w = io.Pipe()

	s.cmd.Stdin = s.r

	if log.IsVerbose() {
		s.cmd.Stdout, s.cmd.Stderr = os.Stdout, os.Stderr
	}

	fmt.Println("cmd", s.cmd)

	return s.cmd.Start()
}

func (s *FFMPegSession) writeBytes(b []byte) (int, error) {
	return s.w.Write(b)
}

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
