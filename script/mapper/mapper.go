package mapper

import (
	"github.com/peter-mount/go-anim/graph"
	"github.com/peter-mount/go-anim/graph/color"
)

// Mapper provides the graphMapper package
type Mapper struct{}

func (_ Mapper) Filter(m graph.Mapper) graph.Filter {
	return m.Filter()
}

func (_ Mapper) MinLevel(r, g, b int) graph.Mapper {
	return color.MinLevel(uint32(r), uint32(g), uint32(b))
}

func (_ Mapper) MaxLevel(r, g, b int) graph.Mapper {
	return color.MaxLevel(uint32(r), uint32(g), uint32(b))
}

func (_ Mapper) Brighten(amount int) graph.Mapper {
	return color.Brighten(uint32(amount))
}

func (_ Mapper) Darken(amount int) graph.Mapper {
	return color.Darken(uint32(amount))
}

func (_ Mapper) Mono() graph.Mapper { return color.Mono }

func (_ Mapper) Mono16() graph.Mapper { return color.Mono16 }

func (_ Mapper) Red() graph.Mapper { return color.Red }

func (_ Mapper) Red16() graph.Mapper { return color.Red16 }

func (_ Mapper) Green() graph.Mapper { return color.Green }

func (_ Mapper) Green16() graph.Mapper { return color.Green16 }

func (_ Mapper) Blue() graph.Mapper { return color.Blue }

func (_ Mapper) Blue16() graph.Mapper { return color.Blue16 }

func (_ Mapper) DeltaRGB(r, g, b int) graph.Mapper {
	d := color.DeltaRGB{
		R: r,
		G: g,
		B: b,
	}
	return d.Apply
}
