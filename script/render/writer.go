package render

import "image"

// Writer exposes the interface provided by FFMPegSession and FrameSession
// to scripts
type Writer interface {
	// Close the Writer, normally used in try-resources
	Close() error
	// WriteBytes the byte slice as a frame
	WriteBytes(b []byte) (int, error)
	// WriteImage writes the image as a frame
	WriteImage(img image.Image) error
	// WriteImageMulti writes the image as num frames
	WriteImageN(img image.Image, num int) error
}
