package time

import (
	"github.com/peter-mount/go-kernel/v2/util"
	"testing"
)

func TestIterator(t *testing.T) {
	type test struct {
		name       string // name of test
		start      string // start timecode
		frameCount int    // if >0 number of frames from start
		until      string // if set then timecode of end
		expect     string // timecode to expect at end
	}

	tests := []test{
		// ====================
		// ForFrames
		// ====================
		{name: "ForFrames 10", start: "00:00:00:00", frameCount: 10, expect: "00:00:00:10"},
		{name: "ForFrames 30", start: "00:00:00:00", frameCount: 30, expect: "00:00:01:00"},
		{name: "ForFrames 60", start: "00:00:00:00", frameCount: 60, expect: "00:00:02:00"},
		{name: "ForFrames 90", start: "00:00:00:00", frameCount: 90, expect: "00:00:03:00"},
		// ====================
		// Until
		// ====================
		{name: "Until 10", start: "00:00:00:00", until: "00:00:00:10", expect: "00:00:00:10"},
		{name: "Until 30", start: "00:00:00:00", until: "00:00:01:00", expect: "00:00:01:00"},
		{name: "Until 60", start: "00:00:00:00", until: "00:00:02:00", expect: "00:00:02:00"},
		{name: "Until 90", start: "00:00:00:00", until: "00:00:03:00", expect: "00:00:03:00"},
		// ====================
		// Cross midnight
		// ====================
		{name: "ForFrames 30 midnight", start: "23:59:59:10", frameCount: 30, expect: "01:00:00:00:10"},
		{name: "ForFrames 60 midnight", start: "23:59:59:20", frameCount: 60, expect: "01:00:00:01:20"},
		{name: "ForFrames 90 midnight", start: "23:59:59:20", frameCount: 90, expect: "01:00:00:02:20"},
		{name: "Until 10 midnight", start: "23:59:59:20", until: "00:00:00:10", expect: "01:00:00:00:10"},
		{name: "Until 30 midnight", start: "23:59:59:20", until: "00:00:01:00", expect: "01:00:00:01:00"},
		{name: "Until 60 midnight", start: "23:59:59:20", until: "00:00:02:00", expect: "01:00:00:02:00"},
		{name: "Until 90 midnight", start: "23:59:59:20", until: "00:00:03:00", expect: "01:00:00:03:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := NewTimeCode(30)
			_, err := start.Set(tt.start)
			if err != nil {
				t.Fatal(err)
				return
			}

			expect, err := ParseTimeCode(tt.expect, 30)
			if err != nil {
				t.Fatal(err)
				return
			}

			var iterator util.Iterator[TimeCodeFragment]
			switch {
			case tt.frameCount > 0:
				iterator = start.ForFrames(tt.frameCount)

			case tt.until != "":
				it, err := start.Until(tt.until)
				if err != nil {
					t.Fatal(err)
					return
				}
				iterator = it
			}

			if iterator == nil {
				t.Fatal("No iterator")
				return
			}

			if !iterator.HasNext() {
				t.Fatal("TimeCodeFragmentIterator empty?")
				return
			}

			var got TimeCodeFragment
			for iterator.HasNext() {
				got = iterator.Next()
			}

			if !expect.Equals(got) {
				t.Errorf("Got %v expected %v", got.TimeCode(), expect.TimeCode())
			}
		})
	}
}
