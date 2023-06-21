package time

import (
	"github.com/peter-mount/go-anim/util"
	"strconv"
	"strings"
)

type Duration struct {
	F float64
	U Unit
}

func (v Duration) IsZero() bool {
	return v.F == 0
}

func (v Duration) Convert(to Unit) Duration {
	if v.U == to {
		return v
	}

	return Duration{
		F: v.F * to.SecondsPer() / v.U.SecondsPer(),
		U: to,
	}
}

// Frames returns the value in frames based on a specific timestamp
func (v Duration) Frames(fps float64) float64 {
	return v.Convert(Second).F * fps
}

func (v Duration) String() string {
	return util.FloatToA(v.F) + v.U.String()
}

func ParseDuration(s string) (Duration, error) {
	unit := Second
	for u, n := range unitNames {
		if strings.HasSuffix(s, n) {
			unit = Unit(u)
			s = strings.TrimSuffix(s, n)
		}
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Duration{}, err
	}

	return Duration{F: f, U: unit}, nil
}

type Unit uint8

const (
	Second Unit = iota
	Minute
	Hour
	Day
)

func (u Unit) String() string {
	return unitNames[u]
}

var unitNames = []string{"s", "m", "h", "d"}

func (u Unit) SecondsPer() float64 {
	switch u {
	case Second:
		return 1.0
	case Minute:
		return 60.0
	case Hour:
		return 3600.0
	case Day:
		return 86400.0
	default:
		return 1.0
	}
}
