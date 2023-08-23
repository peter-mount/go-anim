package clip

import "github.com/peter-mount/go-anim/script/util"

type Frame struct {
	timeCode util.TimeCodeFragment // timeCode of the Frame
	name     string                // Name of frame, e.g. file name if a directory of frames
}
