package time

import (
	"encoding/xml"
	"git.area51.dev/peter/videoident/util"
	util2 "github.com/peter-mount/go-graphics/util"
	"strconv"
)

type Span struct {
	start       int      // Start frame, used in ordering
	end         int      // End frame, used in ordering, duration overrides
	duration    Duration // Duration of the span
	useDuration bool     // true if duration takes precedence over end.
}

// Clear returns a Span which starts where this one does but ends at the same position.
// Duration will be reset but set to the same unit as the source.
func (s Span) Clear() Span {
	return Span{
		start:       s.start,
		end:         s.start,
		duration:    Duration{U: s.duration.U},
		useDuration: s.useDuration,
	}
}

func (s Span) Start() int {
	return s.start
}

func (s Span) End() int {
	return s.end
}

func (s Span) Range() (int, int) {
	return s.start, s.end
}

func (s Span) Length() int {
	return s.end - s.start + 1
}

func (s Span) Duration() Duration {
	return s.duration
}

func (s Span) ContainsFrame(frame int) bool {
	return util.Within(frame, s.start, s.end)
}

func (s Span) Move(ns int) Span {
	l := s.Length()
	s.start = ns
	s.end = ns + l - 1 // -1 to account for start being frame 0
	return s
}

func (s Span) ApplyFPS(fps int) Span {
	if s.useDuration {
		s.end = s.start + (fps * int(s.duration.Convert(Second).F)) - 1
	} else {
		s.duration = s.calcDuration(fps)
	}
	return s
}

func (s Span) calcDuration(fps int) Duration {
	return Duration{F: float64(s.Length()) / float64(fps), U: Second}
}

// Add adds two spans together so the returned one includes both.
func (s Span) Add(b Span, fps int) Span {
	r := Span{
		start: util2.Min(s.start, b.start),
		end:   util2.Max(s.end, b.end),
		// Only keep the duration on the left hand side
		useDuration: s.useDuration,
	}

	// Update
	r.duration = r.calcDuration(fps)

	return r
}

func (s Span) AppendSpanAttr(a []xml.Attr) []xml.Attr {
	//a = util.AppendIntAttr(a, "start", s.start)

	// Include only one of these attributes depending on which one takes precedence
	if s.useDuration {
		a = util.AppendAttr(a, "duration", s.duration.String())
	} else {
		a = util.AppendIntAttr(a, "end", s.end)
	}
	return a
}

func ParseSpan(a []xml.Attr) (Span, error) {
	s := Span{start: -1, end: -1}
	var err error
	for _, attr := range a {
		switch attr.Name.Local {
		//case "start":
		//	s.start, err = strconv.Atoi(attr.Value)

		case "end":
			s.end, err = strconv.Atoi(attr.Value)

		case "duration":
			s.duration, err = ParseDuration(attr.Value)
			s.useDuration = true
		}

		if err != nil {
			return Span{}, err
		}
	}

	return s, nil
}
