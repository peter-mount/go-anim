package time

import "testing"

func TestTimeCodeFragment(t *testing.T) {
	type test struct {
		name      string // test name
		left      string // left timecode
		right     string // right timecode
		addFrames int    // if != 0 then add to left then compare with right
		equals    bool   // true if Equals() must return true
		before    bool   // true if Before() must return true
		after     bool   // true if After() must return true
		leftDay   int    // day value of left
		rightDay  int    // day value of right
	}

	tests := []test{
		{name: "zero", left: "00:00:00", right: "00:00:00", equals: true},
		// ========================================
		// Basic comparisons at the frame level
		// ========================================
		{name: "frame equals", left: "00:00:00:01", right: "00:00:00:01", equals: true},
		{name: "frame before", left: "00:00:00:00", right: "00:00:00:01", before: true},
		{name: "frame after", left: "00:00:00:01", right: "00:00:00:00", after: true},
		// ========================================
		// Basic comparisons at the second level
		// ========================================
		{name: "second equals", left: "00:00:01", right: "00:00:01", equals: true},
		{name: "second before", left: "00:00:00", right: "00:00:01", before: true},
		{name: "second after", left: "00:00:01", right: "00:00:00", after: true},
		// ========================================
		// Basic comparisons at the minute level
		// ========================================
		{name: "minute equals", left: "00:01:00", right: "00:01:00", equals: true},
		{name: "minute before", left: "00:00:00", right: "00:01:00", before: true},
		{name: "minute after", left: "00:01:00", right: "00:00:00", after: true},
		// ========================================
		// Basic comparisons at the hour level
		// ========================================
		{name: "hour equals", left: "01:00:00", right: "01:00:00", equals: true},
		{name: "hour before", left: "01:00:00", right: "02:00:00", before: true},
		{name: "hour after", left: "02:00:00", right: "01:00:00", after: true},
		// ========================================
		// Comparisons crossing midnight
		// ========================================
		{name: "midnight equals", left: "01:00:00", right: "01:00:00", leftDay: 1, rightDay: 1, equals: true},
		{name: "hour before", left: "01:00:00", right: "02:00:00", leftDay: 1, rightDay: 1, before: true},
		{name: "hour after", left: "02:00:00", right: "01:00:00", leftDay: 1, rightDay: 1, after: true},
		// ========================================
		// Add frames
		// ========================================
		{name: "AddFrames 10", left: "00:00:00", addFrames: 10, right: "00:00:00:10"},
		{name: "AddFrames 60", left: "00:00:00", addFrames: 60, right: "00:00:02:00"},
		// test adding crosses midnight correctly
		{name: "AddFrames 60 midnight", left: "23:59:59", addFrames: 60, right: "00:00:01:00", rightDay: 1},
	}

	testFragment := func(t *testing.T, a test, want bool, f func(TimeCodeFragment, TimeCodeFragment) bool) {
		left, err := ParseTimeCode(a.left, 30)
		if err != nil {
			t.Fatal(err)
			return
		}
		left.day = a.leftDay

		right, err := ParseTimeCode(a.right, 30)
		if err != nil {
			t.Fatal(err)
			return
		}
		right.day = a.rightDay

		switch {
		case a.addFrames != 0:
			result := left.AddFrames(a.addFrames)
			if got := result.Equals(right); got != want {
				t.Errorf("equals got = %v, want %v", got, want)
			}

		default:
			if got := f(left, right); got != want {
				t.Errorf("got = %v, want %v", got, want)
			}
		}
	}

	for _, tt := range tests {
		switch {
		case tt.addFrames != 0:
			t.Run(tt.name, func(t *testing.T) {
				testFragment(t, tt, true, nil)
			})

		default:
			t.Run(tt.name+" Equals", func(t *testing.T) {
				testFragment(t, tt, tt.equals, TimeCodeFragment.Equals)
			})
			t.Run(tt.name+" Before", func(t *testing.T) {
				testFragment(t, tt, tt.before, TimeCodeFragment.Before)
			})
			t.Run(tt.name+" After", func(t *testing.T) {
				testFragment(t, tt, tt.after, TimeCodeFragment.After)
			})
		}
	}
}
