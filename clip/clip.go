package clip

import "github.com/peter-mount/go-anim/script/util"

// Clip represents a collection of Frame's that represent a sequence to be rendered into an animation.
type Clip struct {
	timeCode util.TimeCode // TimeCode of the entire Clip
	frames   []Frame
}
