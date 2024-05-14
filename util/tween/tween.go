package tween

import (
	"fmt"
	"github.com/peter-mount/go-anim/util"
	"math"
)

// Tween (short for between) provides a basic form of animation control.
type Tween interface {
	// ContainsFrame returns true if the specified frame is within this Tween
	ContainsFrame(int) bool
	// T for the specified frame.
	// The result will be either a value between 0 and 1 inclusive.
	//
	// T==0 when frame==start,
	// T==1 when frame==end,
	// T==NaN if frame is not within start...end
	T(int) float64
	// Start frame
	Start() int
	// End frame
	End() int
	// FrameLength of the Tween in frames.
	// The returned value will always be greater than 0.
	// This is the same as End()-Start()+1
	FrameLength() int
}

// T calculates a value between 0 and 1 based on the supplied frame.
// This will only be called if Tween.ContainsFrame(frame) is true.
// Here, frame is the frame offset from the start of the tween, so
// frame ranges from 0 to Tween.FrameLength inclusive.
type T func(frame int) float64

// TFactory is a function that will return T based on a specific Tween.
type TFactory func(Tween) T

// New creates a new Tween using a TFactory which applies for a specific frame range.
// The Tween generated will be valid from startFrame to endFrame inclusively.
// If factory is nil then this will panic.
// Internally it ensures that startFrame is before endFrame.
func New(startFrame, endFrame int, factory TFactory) Tween {
	if startFrame > endFrame {
		startFrame, endFrame = endFrame, startFrame
	}

	t := &basicTween{
		startFrame:  startFrame,
		endFrame:    endFrame,
		frameLength: endFrame - startFrame + 1,
	}

	t.t = factory(t)

	return t
}

// NewDuration returns a frame which starts at startFrame and runs for duration frames.
// This will panic if duration < 1 or factory is nil
func NewDuration(startFrame, duration int, factory TFactory) Tween {
	if duration < 1 {
		panic(fmt.Errorf("invalid duration %d", duration))
	}

	return New(startFrame, startFrame+duration-1, factory)
}

type basicTween struct {
	startFrame  int
	endFrame    int
	frameLength int
	t           T
}

func (t *basicTween) ContainsFrame(frame int) bool {
	return util.Within(frame, t.startFrame, t.endFrame)
}

func (t *basicTween) T(frame int) float64 {
	if !t.ContainsFrame(frame) {
		return math.NaN()
	}
	return t.t(frame - t.startFrame)
}

func (t *basicTween) Start() int {
	return t.startFrame
}

func (t *basicTween) End() int {
	return t.endFrame
}

func (t *basicTween) FrameLength() int {
	return t.frameLength
}
