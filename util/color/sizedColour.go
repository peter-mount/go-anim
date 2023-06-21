package color

import (
	"fmt"
	"github.com/peter-mount/go-anim/util/unit"
	"image/color"
	"strings"
)

// SizedColor is a size & colour but represented as a single xml attribute
type SizedColor struct {
	Size  unit.Value
	Color color.Color
}

func (s SizedColor) String() string {
	var a []string
	if !s.Size.IsZero() {
		a = append(a, s.Size.String())
	}
	if s.Color != nil {
		a = append(a, ColourString(s.Color))
	}
	return strings.Join(a, " ")
}

func ParseSizedColor(s string) (SizedColor, error) {
	var r SizedColor
	var err error
	a := strings.Split(s, " ")
	switch len(a) {
	case 0:
		return r, nil

	case 1:
		r.Size, err = unit.ParseValue(a[0])
		if err != nil {
			r.Color, err = ParseColour(a[0])
		}
		return r, err

	case 2:
		r.Size, err = unit.ParseValue(a[0])
		if err == nil {
			r.Color, err = ParseColour(a[1])
		}
		return r, err

	default:
		return r, fmt.Errorf("invalid SizedColor %q", s)
	}
}
