package tween

type Frame interface {
	// Frame is the absolute frame in the animation, starts from 1
	Frame() int
	// Offset is the frame number within the Tween, starts from 0
	Offset() int
	// T is the fraction between the Tween's start and end frames,
	// 0 when Frame() is the start frame, 1 when Frame() is the end frame.
	T() float64
	// Tween is the Tween this Frame is for
	Tween() Tween
}

func NewFrame(t Tween, frame int) Frame {
	return &basicFrame{
		tween:  t,
		frame:  frame,
		offset: frame - t.Start(),
		t:      t.T(frame),
	}
}

type basicFrame struct {
	tween  Tween
	frame  int
	offset int
	t      float64
}

func (f *basicFrame) Frame() int {
	return f.frame
}

func (f *basicFrame) Offset() int {
	return f.offset
}

func (f *basicFrame) T() float64 {
	return f.t
}

func (f *basicFrame) Tween() Tween {
	return f.tween
}
