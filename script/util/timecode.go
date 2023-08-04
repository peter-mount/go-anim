package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TimeCode handles the management of Timecodes
type TimeCode struct {
	frameNum   int // The overall frame number
	startSec   int // The time code in seconds since 00:00:00:ff of the start
	startFrame int // The initial frame number at the start, usually 0
	sec        int // The current time code
	frameSec   int // Frame within the current second
	frameRate  int // Frame Rate
}

func NewTimeCode(frameRate int) *TimeCode {
	return &TimeCode{frameNum: 1, frameRate: frameRate}
}

// FrameNum is the overall frame number, starting at 1.
// This can be used when forming file names for individual frame images
func (tc *TimeCode) FrameNum() int {
	return tc.frameNum
}

// Sec returns the current whole second within the clip since 00:00:00
func (tc *TimeCode) Sec() int {
	return tc.sec
}

// StartSec returns the starting whole second since 00:00:00 for the clip.
func (tc *TimeCode) StartSec() int {
	return tc.startSec
}

// FrameSec returns the frame within the current second
func (tc *TimeCode) FrameSec() int {
	return tc.frameSec
}

// FrameRate returns the frame rate of the clip
func (tc *TimeCode) FrameRate() int {
	return tc.frameRate
}

// StartTimeCode returns the start time code as "hh:mm:ss:ff"
func (tc *TimeCode) StartTimeCode() string {
	return tc.timeCode(tc.startSec, tc.startFrame)
}

// TimeCode returns the current time code as "hh:mm:ss:ff"
func (tc *TimeCode) TimeCode() string {
	return tc.timeCode(tc.sec, tc.frameSec)
}

func (tc *TimeCode) timeCode(sec, frame int) string {
	m := sec / 60
	h := m / 60
	return fmt.Sprintf("%02d:%02d:%02d:%02d", h%24, m%60, sec%60, frame%tc.frameRate)
}

// Next moves the TimeCode to the next frame.
func (tc *TimeCode) Next() {
	tc.frameSec++
	if tc.frameSec >= tc.frameRate {
		tc.frameSec = 0
		tc.sec++

		if tc.sec >= 86400 {
			tc.sec = 0
		}
	}
}

// Set sets the starting TimeCode. This is in the format "hh:mm:ss:ff" although a short form "hh:mm:ss" is
// valid in which case ff will be 0.
//
// This will return an error if the TimeCode has been used for a frame, e.g. Next() has been called.
func (tc *TimeCode) Set(s string) (*TimeCode, error) {
	if tc.frameNum > 1 {
		return nil, errors.New("cannot Set a running TimeCode")
	}
	a := strings.Split(s, ":")

	l := len(a)

	// Allow either "hh:mm:ss" or "hh:mm:ss:ff". For the shorter, ff=0
	valid := l == 3 || l == 4

	var v []int
	for i := 0; i < l && valid; i++ {
		n, err := strconv.Atoi(a[i])
		valid = err == nil && n >= 0
		switch i {
		case 0:
			valid = valid && n < 24
		case 1, 2:
			valid = valid && n < 60
		case 3:
			valid = valid && n < tc.frameRate
		}
		if valid {
			v = append(v, n)
		}
	}

	if !valid {
		return nil, fmt.Errorf("invalid timecode %q must be hh:mm:ss or hh:mm:ss:00", s)
	}

	// If short form, set start frame to 0
	if l == 3 {
		v = append(v, 0)
	}

	tc.startSec = (((v[0] * 60) + v[1]) * 60) + v[2]
	tc.startFrame = v[3]

	tc.sec, tc.frameSec = tc.startSec, tc.startFrame

	return tc, nil
}

// IsRunning returns true if the TimeCode is running. Specifically once a frame has been rendered it is running.
func (tc *TimeCode) IsRunning() bool {
	return tc.frameNum > 1
}
