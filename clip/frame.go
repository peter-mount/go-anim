package clip

import (
	"github.com/peter-mount/go-anim/util/time"
)

type Frame struct {
	timeCode time.TimeCodeFragment // timeCode of the Frame
	name     string                // Name of frame, e.g. file name if a directory of frames
}
