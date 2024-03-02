package tween

// Linear is a TFactory which creates a linear tween where the point moves at a constant rate
// between the start and end frames.
func Linear(tween Tween) T {
	delta := 1.0 / float64(tween.FrameLength())

	return func(frame int) float64 {
		return float64(frame) * delta
	}
}
