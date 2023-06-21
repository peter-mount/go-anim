package util

// MinMax returns the min and max values for a slice
func MinMax(a []uint32) (uint32, uint32) {
	var min, max uint32
	for i, v := range a {
		if i == 0 {
			min = v
			max = v
		} else {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	return min, max
}

// MaxSlice returns the max values for a slice
func MaxSlice(a []uint32) uint32 {
	var max uint32
	for i, v := range a {
		if i == 0 {
			max = v
		} else if v > max {
			max = v
		}
	}
	return max
}

func MaxU8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func MinU8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func MaxU32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func MinU32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// Limit limits d so it lies within a and b inclusively.
func Limit(d, a, b int) int {
	if d < a {
		return a
	}
	if d > b {
		return b
	}
	return d
}

// LimitU8 limits d within 0...255
func LimitU8(d int) uint8 {
	return uint8(Limit(d, 0, 255))
}

// LimitU16 limits d within 0...65535
func LimitU16(d int) uint16 {
	return uint16(Limit(d, 0, 65535))
}

// LimitU32 is identical to LimitU16 but used as image.RGBA uses uint32 with uint16 limits
func LimitU32(d int) uint32 {
	return uint32(LimitU16(d))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Within returns true if frame is within start...end inclusively
func Within(frame, start, end int) bool {
	return frame >= start && frame <= end
}
