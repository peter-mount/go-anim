package color

import (
	"github.com/peter-mount/go-anim/graph"
	"image"
	"image/color"
)

// Histogram represents the total number of pixels in an image based on RGB values
type Histogram struct {
	Red   []uint32
	Green []uint32
	Blue  []uint32
}

const (
	maxIndex = 255
)

func NewHistogram() *Histogram {
	return &Histogram{
		Red:   make([]uint32, maxIndex+1),
		Green: make([]uint32, maxIndex+1),
		Blue:  make([]uint32, maxIndex+1),
	}
}

// Add adds a color.Color to the Histogram
func (h *Histogram) Add(c color.Color) *Histogram {
	r, g, b, _ := c.RGBA()
	h.Red[r>>8]++
	h.Green[g>>8]++
	h.Blue[b>>8]++
	return h
}

// AnalyzeImage analyzes an image placing the results in the histogram
func (h *Histogram) AnalyzeImage(src image.Image) *Histogram {
	if src != nil {
		_ = graph.Of(h.AnalyzeFilter).
			Do(src, nil, src.Bounds())
	}
	return h
}

// AnalyzeFilter analyzes an image placing the results in the histogram
func (h *Histogram) AnalyzeFilter(_, _ int, c color.Color) (color.Color, error) {
	h.Add(c)
	return c, nil
}

// ResetValue sets the histogram values at x to 0.
// Used after Analyze to remove high known values which affects the results.
// Eg an image of mostly black needs black removed.
func (h *Histogram) ResetValue(x int) *Histogram {
	h.Red[x] = 0
	h.Green[x] = 0
	h.Blue[x] = 0
	return h
}

// ResetValuesBelow sets the histogram values below v to 0.
// Used after Analyze to remove high known values which affects the results.
// Eg an image of mostly black needs black removed.
func (h *Histogram) ResetValuesBelow(v int) *Histogram {
	for x := 0; x <= v; x++ {
		h.ResetValue(x)
	}
	return h
}

// ResetValuesAbove sets the histogram values above v to 0.
// Used after Analyze to remove high known values which affects the results.
// Eg an image of mostly white needs white
func (h *Histogram) ResetValuesAbove(v int) *Histogram {
	for x := v; x <= maxIndex; x++ {
		h.ResetValue(x)
	}
	return h
}
