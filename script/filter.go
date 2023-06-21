package script

import (
	"github.com/peter-mount/go-anim/graph/color"
	"github.com/peter-mount/go-anim/graph/filter"
	"image"
)

import (
	"github.com/peter-mount/go-anim/graph"
)

type Filter struct{}

// Filter applies a graph.Filter on a source image within the specified bounds,
// writing the result to the destination image.
//
// The source and destination image may be the same Image if the filter supports it.
func (_ Filter) Filter(f graph.Filter, src image.Image, dst graph.Image, b image.Rectangle) error {
	return f.Do(src, dst, b)
}

// FilterNew applies a graph.Filter on a source image,
// returning a new mutable image with the new content.
func (_ Filter) FilterNew(f graph.Filter, src image.Image) (graph.Image, error) {
	return f.DoNew(src)
}

// FilterOver applies the filter over the supplied mutable image,
// overwriting its previous state.
func (_ Filter) FilterOver(f graph.Filter, src graph.Image) error {
	return f.DoOver(src)
}

func (_ Filter) Equalize(h *color.Histogram, b image.Rectangle) graph.Filter {
	return filter.EqualizeFilter(h, b)
}
