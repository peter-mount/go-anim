package time

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TimeCode handles the management of Timecodes during an animation
type TimeCode struct {
	frameNum int              // The overall frame number
	start    TimeCodeFragment // Start TimeCode
	current  TimeCodeFragment // TimeCode of the current frame
}

const (
	// The initial frameNum = 1 as it's the index used in filenames
	tcStartFrameNum = 1
)

// NewTimeCode creates a TimeCode with the specified frame rate
func NewTimeCode(frameRate int) *TimeCode {
	return &TimeCode{
		frameNum: tcStartFrameNum,
		start:    TimeCodeFragment{frameRate: frameRate},
		current:  TimeCodeFragment{frameRate: frameRate},
	}
}

// FrameRate returns the frame rate of the clip
func (tc *TimeCode) FrameRate() int {
	return tc.start.FrameRate()
}

func (tc *TimeCode) FrameRateF() float64 {
	return tc.start.FrameRateF()
}

// FrameNum is the overall frame number, starting at 1.
// This can be used when forming file names for individual frame images
func (tc *TimeCode) FrameNum() int {
	return tc.frameNum
}

// StartTimeCode returns the time code of the first frame
func (tc *TimeCode) StartTimeCode() TimeCodeFragment {
	return tc.start
}

// TimeCode returns the time code for the current frame
func (tc *TimeCode) TimeCode() TimeCodeFragment {
	return tc.current
}

// Next moves the TimeCode to the next frame.
func (tc *TimeCode) Next() {
	// Next frame serial
	tc.frameNum++

	// Don't optimise caching tc.current as it's not a pointer
	tc.current.frame++
	if tc.current.frame >= tc.current.frameRate {
		tc.current.frame = 0
		tc.current.sec++

		if tc.current.sec >= 86400 {
			tc.current.sec = 0
			tc.current.day++
		}
	}
}

// Set sets the starting TimeCode. This is in the format "hh:mm:ss:ff" although a short form "hh:mm:ss" is
// valid in which case ff will be 0.
//
// This will return an error if the TimeCode has been used for a frame, e.g. Next() has been called.
func (tc *TimeCode) Set(s string) (*TimeCode, error) {
	if tc.IsRunning() {
		return nil, errors.New("cannot Set a running TimeCode")
	}

	tcf, err := ParseTimeCode(s, tc.FrameRate())
	if err != nil {
		return nil, err
	}

	tc.start = tcf
	tc.current = tc.start

	return tc, nil
}

// IsRunning returns true if the TimeCode is running. Specifically once a frame has been rendered it is running.
func (tc *TimeCode) IsRunning() bool {
	return tc.frameNum > tcStartFrameNum
}

type TimeCodeFragment struct {
	day       int // The day, used in handling fragments spanning midnight
	sec       int // The current time code
	frame     int // Frame within the current second
	frameRate int // Frame Rate
}

func ParseTimeCode(s string, frameRate int) (TimeCodeFragment, error) {
	a := strings.Split(s, ":")

	l := len(a)

	// Allow either "hh:mm:ss" or "hh:mm:ss:ff". For the shorter, ff=0
	valid := l == 3 || l == 4

	// "dd:hh:mm:ss:ff" allows for spanning midnight
	day := 0
	if l == 5 {
		n, err := strconv.Atoi(a[0])
		valid = err == nil && n >= 0
		day = n
		// Strip off the day from the slice
		a = a[1:]
		l--
	}

	// Parse each field as an int, testing for bounds
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
			valid = valid && n < frameRate
		}
		if valid {
			v = append(v, n)
		}
	}

	tc := TimeCodeFragment{frameRate: frameRate, day: day}

	if !valid {
		return tc, fmt.Errorf("invalid timecode %q must be hh:mm:ss or hh:mm:ss:ff", s)
	}

	// If short form, set start frame to 0
	if l == 3 {
		v = append(v, 0)
	}

	tc.sec = (((v[0] * 60) + v[1]) * 60) + v[2]
	tc.frame = v[3]

	return tc, nil
}

func (tc TimeCodeFragment) TimeCode() string {
	return fmt.Sprintf("%02d:%02d:%02d:%02d", tc.Hour(), tc.Minute(), tc.Second(), tc.Frame())
}

// FrameRate of the clip
func (tc TimeCodeFragment) FrameRate() int {
	return tc.frameRate
}

func (tc TimeCodeFragment) FrameRateF() float64 {
	return float64(tc.Frame())
}

// Offset returns the number of seconds since "00:00:00" for the clip
func (tc TimeCodeFragment) Offset() int {
	return tc.sec + (tc.day * 86400)
}

func (tc TimeCodeFragment) Day() int {
	return tc.day
}

// Hour returns the hour component as an int
func (tc TimeCodeFragment) Hour() int {
	return (tc.sec / 3600) % 24
}

// Minute returns the minute component as an int
func (tc TimeCodeFragment) Minute() int {
	return (tc.sec / 60) % 60
}

// Second returns the second component as an int
func (tc TimeCodeFragment) Second() int {
	return tc.sec % 60
}

func (tc TimeCodeFragment) DayF() float64 {
	return float64(tc.Day())
}

func (tc TimeCodeFragment) HourF() float64 {
	return float64(tc.Hour())
}

func (tc TimeCodeFragment) MinuteF() float64 {
	return float64(tc.Minute())
}

func (tc TimeCodeFragment) SecondF() float64 {
	return float64(tc.Second())
}

// DayS returns the hour component as a 2 digit string, useful in rendering
func (tc TimeCodeFragment) DayS() string {
	return tc.digit(tc.Day())
}

// HourS returns the hour component as a 2 digit string, useful in rendering
func (tc TimeCodeFragment) HourS() string {
	return tc.digit(tc.Hour())
}

// MinuteS returns the minute component as a 2 digit string, useful in rendering
func (tc TimeCodeFragment) MinuteS() string {
	return tc.digit(tc.Minute())
}

// SecondS returns the second component as a 2 digit string, useful in rendering
func (tc TimeCodeFragment) SecondS() string {
	return tc.digit(tc.Second())
}

// Frame returns the frame within the current second
func (tc TimeCodeFragment) Frame() int {
	return tc.frame
}

func (tc TimeCodeFragment) FrameF() float64 {
	return float64(tc.Frame())
}

// FrameS returns the frame component as a 2 digit string, useful in rendering
func (tc TimeCodeFragment) FrameS() string {
	return tc.digit(tc.Frame())
}

func (tc TimeCodeFragment) digit(n int) string {
	s := strconv.Itoa(n)
	if len(s) == 1 {
		s = "0" + s
	}
	return s
}

// FramesRemaining returns the number of frames remaining in the current second.
// This is simply frameRate - frame.
func (tc TimeCodeFragment) FramesRemaining() int {
	return tc.FrameRate() - tc.Frame()
}

// IsStartSecond returns true if the current frame is the first of a whole second.
func (tc TimeCodeFragment) IsStartSecond() bool {
	return tc.Frame() == 0
}

func (tc *TimeCode) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, tc.frameNum)
	if err == nil {
		err = tc.start.Write(w)
	}
	return err
}

func ReadTimeCode(r io.Reader) (TimeCode, error) {
	var tc TimeCode
	err := binary.Read(r, binary.BigEndian, &tc.frameNum)
	if err == nil {
		tc.start, err = ReadTimeCodeFragment(r)
	}
	if err == nil {
		tc.current = tc.start
	}
	return tc, err
}

func (tc TimeCodeFragment) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, tc.sec)
	if err == nil {
		err = binary.Write(w, binary.BigEndian, tc.frame)
	}
	if err == nil {
		err = binary.Write(w, binary.BigEndian, tc.frameRate)
	}
	return err
}

func ReadTimeCodeFragment(r io.Reader) (TimeCodeFragment, error) {
	var tc TimeCodeFragment
	err := binary.Read(r, binary.BigEndian, &tc.sec)
	if err == nil {
		err = binary.Read(r, binary.BigEndian, &tc.frame)
	}
	if err == nil {
		err = binary.Read(r, binary.BigEndian, &tc.frameRate)
	}
	return tc, err
}

func (tc TimeCodeFragment) Equals(b TimeCodeFragment) bool {
	return tc.day == b.day && tc.sec == b.sec && tc.frame == b.frame
}

func (tc TimeCodeFragment) Before(b TimeCodeFragment) bool {
	if tc.day != b.day {
		return tc.day < b.day
	}

	return tc.sec < b.sec || (tc.sec == b.sec && tc.frame < b.frame)
}

func (tc TimeCodeFragment) NotAfter(b TimeCodeFragment) bool {
	return tc.Before(b) || tc.Equals(b)
}

func (tc TimeCodeFragment) After(b TimeCodeFragment) bool {
	if tc.day != b.day {
		return tc.day > b.day
	}

	return tc.sec > b.sec || (tc.sec == b.sec && tc.frame > b.frame)
}

func (tc TimeCodeFragment) toFrame(d, h, m, s, f int) int {
	return (((d * 86400) + (h * 3600) + (m * 60) + s) * tc.frameRate) + f
}

func (tc TimeCodeFragment) AddFrames(count int) TimeCodeFragment {
	if count < 1 {
		return tc
	}

	return tc.new(tc.toFrame(tc.day, 0, 0, tc.sec, tc.frame) + count)
}

func (tc TimeCodeFragment) new(f int) TimeCodeFragment {
	sec := f / tc.frameRate

	return TimeCodeFragment{
		day:       sec / 86400,
		sec:       sec % 86400,
		frame:     f % tc.frameRate,
		frameRate: tc.frameRate,
	}
}

func (tc TimeCodeFragment) Add(d, h, m, s, f int) TimeCodeFragment {
	return tc.AddFrames(tc.toFrame(d, h, m, s, f))
}
