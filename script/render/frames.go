package render

import (
	"fmt"
	"github.com/peter-mount/go-anim/util/time"
	"os"
)

// FrameSession handles sending frames to individual images on disk
type FrameSession struct {
	RenderStreamBase
}

func (_ Render) frames(fileName string, frameRate int, encoder Encoder) *FrameSession {
	s := &FrameSession{
		RenderStreamBase: RenderStreamBase{
			fileName: fileName,
			timeCode: time.NewTimeCode(frameRate),
			encoder:  encoder,
		},
	}
	s.RenderStreamBase.write = s.writeBytes
	return s
}

func (s *FrameSession) writeBytes(b []byte) (int, error) {
	fileName := fmt.Sprintf(s.fileName, s.TimeCode().FrameNum())

	return 0, os.WriteFile(fileName, b, 0644)
}

// Close closes the FrameSession.
// Normally this is handled by try-resources
func (s *FrameSession) Close() error {
	// For now does nothing, just to implement the Writer interface
	return nil
}
