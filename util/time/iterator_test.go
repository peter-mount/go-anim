package time

import (
	"testing"
)

func TestIterator(t *testing.T) {
	type test struct {
		name       string // name of test
		start      string // start timecode
		frameCount int    // if >0 number of frames from start
		until      string // if set then timecode of end
		expect     string // timecode to expect at end
		expectDay  int    // day of end timecode
	}

	tests := []test{
		// ====================
		// Test from 0
		// ====================
		{name: "zero 10", start: "00:00:00:00", frameCount: 10, expect: "00:00:00:10"},
		{name: "zero 30", start: "00:00:00:00", frameCount: 30, expect: "00:00:01:00"},
		{name: "zero 60", start: "00:00:00:00", frameCount: 60, expect: "00:00:02:00"},
		{name: "zero 90", start: "00:00:00:00", frameCount: 90, expect: "00:00:03:00"},
		{name: "until 10", start: "00:00:00:00", until: "00:00:00:10", expect: "00:00:00:10"},
		{name: "until 30", start: "00:00:00:00", until: "00:00:01:00", expect: "00:00:01:00"},
		{name: "until 60", start: "00:00:00:00", until: "00:00:02:00", expect: "00:00:02:00"},
		{name: "until 90", start: "00:00:00:00", until: "00:00:03:00", expect: "00:00:03:00"},
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
			expect.day = tt.expectDay

			var iterator *Iterator
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
				t.Fatal("Iterator empty?")
				return
			}

			var got TimeCodeFragment
			for iterator.HasNext() {
				v := iterator.Next()
				if tcf, ok := v.(TimeCodeFragment); ok {
					got = tcf
				}
			}

			if !expect.Equals(got) {
				t.Errorf("Got %v expected %v", got.TimeCode(), expect.TimeCode())
			}
		})
	}
}
