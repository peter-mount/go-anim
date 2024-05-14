package tween

import (
	"github.com/peter-mount/go-anim/util"
	"math"
)

type CompoundTween interface {
	Tween
	// Add appends a new Tween which runs for a specific duration
	Add(duration int, factory TFactory) CompoundTween
}

type compound struct {
	startFrame  int
	endFrame    int
	frameLength int
	tweens      []Tween
}

func NewCompound(startFrame int) Tween {
	return &compound{
		startFrame: startFrame,
		endFrame:   startFrame,
	}
}

func (t *compound) Add(duration int, factory TFactory) CompoundTween {
	var startFrame int
	if len(t.tweens) == 0 {
		startFrame = t.startFrame
	} else {
		startFrame = t.endFrame + 1
	}

	tween := NewDuration(startFrame, duration, factory)
	t.tweens = append(t.tweens, tween)

	t.endFrame = tween.End()
	t.frameLength = t.endFrame - t.startFrame + 1

	return t
}

func (t *compound) ContainsFrame(frame int) bool {
	return t != nil && util.Within(frame, t.startFrame, t.endFrame)
}

func (t *compound) Start() int {
	return t.startFrame
}

func (t *compound) End() int {
	return t.endFrame
}

func (t *compound) FrameLength() int {
	return t.frameLength
}

func (t *compound) locateTween(frame int) Tween {
	for _, c := range t.tweens {
		if c.ContainsFrame(frame) {
			return c
		}
	}
	return nil
}

func (t *compound) T(frame int) float64 {
	var c Tween

	if t.ContainsFrame(frame) {
		c = t.locateTween(frame)
	}

	if c != nil {
		return c.T(frame)
	}

	return math.NaN()
}
