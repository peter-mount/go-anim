package frames

import (
	time2 "github.com/peter-mount/go-anim/util/time"
	"image"
	"time"
)

type Frame struct {
	Source   string    // Filename being processed
	Time     time.Time // Time of frame being processed
	previous *Frame
	image    image.Image
}

func (f *Frame) getRoot() *Frame {
	cf := f
	for cf.previous != nil {
		cf = cf.previous
	}
	return cf
}

func (f *Frame) RequiresImage() bool {
	return f.getRoot().image == nil
}

func (f *Frame) Image() image.Image {
	return f.getRoot().image
}

func (f *Frame) SetImage(img image.Image) {
	f.getRoot().image = img
}

type FrameSet struct {
	frames []*Frame
}

func (fs *FrameSet) Size() int {
	return len(fs.frames)
}

func (fs *FrameSet) HasNext() bool {
	return len(fs.frames) > 0
}

func (fs *FrameSet) Next() *Frame {
	f := fs.frames[0]

	// This enables the frame to be released when it's no longer referenced by future frames sharing its Image reference.
	// Without this, we would run out of memory holding image copies longer than necessary.
	fs.frames[0] = nil

	// Now reduce the slice
	fs.frames = fs.frames[1:]

	return f
}

// Sequence returns a FrameSet containing the frames to be rendered.
// This is the same as SequenceIn except it uses the local TimeZone.
//
// interval is the period in seconds between frames in sourceFiles.
// If frames are missing or are separated by more than this interval then
// this will repeat the previous frames to fill in the gaps.
//
// sourceFiles is a slice of image filenames, with the names being their timestamp.
func Sequence(interval int, sourceFiles []string) *FrameSet {
	return SequenceIn(interval, sourceFiles, time.Local)
}

// SequenceIn returns a FrameSet containing the frames to be rendered.
//
// interval is the period in seconds between frames in sourceFiles.
// If frames are missing or are separated by more than this interval then
// this will repeat the previous frames to fill in the gaps.
//
// sourceFiles is a slice of image filenames, with the names being their timestamp.
//
// loc is the timezone of the images
func SequenceIn(interval int, sourceFiles []string, loc *time.Location) *FrameSet {
	// Step between each frame
	step := time.Duration(interval) * time.Second
	// this is 1.5x step, used when checking if we need an intermediate step.
	// If we used step and the next frame is a second late then we get a step which
	// isn't needed, so this ensures we have a buffer between them
	step2 := time.Duration(float64(interval)*1.5) * time.Second

	frames := &FrameSet{}

	var lastFrame *Frame
	for _, sourceFile := range sourceFiles {
		frame := &Frame{
			Source: sourceFile,
			Time:   time2.TimeFromFileNameIn(sourceFile, loc),
		}

		// Ensure we have frames between steps
		if lastFrame != nil && step > 0 {
			for frame.Time.Sub(lastFrame.Time) >= step2 {
				lastFrame = &Frame{
					Source:   lastFrame.Source,
					Time:     lastFrame.Time.Add(step),
					previous: lastFrame,
				}
				frames.frames = append(frames.frames, lastFrame)
			}
		}

		frames.frames = append(frames.frames, frame)
		lastFrame = frame
	}

	return frames
}
