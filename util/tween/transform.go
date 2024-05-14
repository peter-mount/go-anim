package tween

// Transform
type Transform interface {
	Translate(frameLength int, t float64) Point
}

func NewMotion(start, end Point) Transform {
	return Line{
		Start: start,
		End:   end,
		delta: end.Sub(start),
	}
}

type motion struct {
	start, end Point
	delta      Point
}

func (m motion) Translate(frameLength int, t float64) Point {
	d := t / float64(frameLength)
	return m.start.Translate(m.delta.X*d, m.delta.Y*d)
}
