package io

import (
	"archive/tar"
	"fmt"
	time2 "github.com/peter-mount/go-anim/util/time"
	"image"
	"io"
	"os"
	"time"
)

// TarWriter writes frames to a tar file rather than a directory.
type TarWriter struct {
	RenderStreamBase
	ext string         // Frame extension
	fw  io.WriteCloser // File handle
	tw  *tar.Writer    // Tar writer
}

func (_ Render) tar(fileName string, frameRate int, encoder Encoder, ext string) RenderStream {
	s := &TarWriter{
		RenderStreamBase: RenderStreamBase{
			fileName: fileName,
			timeCode: time2.NewTimeCode(frameRate),
			encoder:  encoder,
		},
		ext: ext,
	}
	s.RenderStreamBase.init = s.init
	s.RenderStreamBase.write = s.writeBytes
	return s
}

func (s *TarWriter) init(_ image.Image) error {
	if s.fw != nil {
		return nil
	}

	f, err := os.Create(s.fileName)
	if err != nil {
		return err
	}

	s.fw = f
	s.tw = tar.NewWriter(f)
	return nil
}

func (s *TarWriter) Close() (err error) {
	if s.tw != nil {
		err = s.tw.Flush()
		err1 := s.tw.Close()
		if err == nil {
			err = err1
		}
		s.tw = nil
	}

	if s.fw != nil {
		err1 := s.fw.Close()
		if err == nil {
			err = err1
		}
	}

	return err
}

func (s *TarWriter) writeBytes(b []byte) (int, error) {
	now := time.Now()

	header := &tar.Header{
		// 8 digit frame number will allow frame rates of 120/s or 1000/s
		// without overflowing within a day's images.
		// 60/s for a day is only 1 digit shorter
		Name:       fmt.Sprintf("%08d%s", s.TimeCode().FrameNum(), s.ext),
		Mode:       0600,
		Size:       int64(len(b)),
		ModTime:    now, // TODO take from TimeCode?
		AccessTime: now,
		ChangeTime: now,
	}

	err := s.tw.WriteHeader(header)
	if err != nil {
		return 0, err
	}

	return s.tw.Write(b)
}
