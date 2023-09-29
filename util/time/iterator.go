package time

// Iterator is returned by TimeCode for advancing it for a set period.
// Each call to Next() will advance the TimeCode.
type Iterator struct {
	tc  *TimeCode
	end TimeCodeFragment
}

func (i *Iterator) HasNext() bool {
	return i.tc.TimeCode().Before(i.end)
}

func (i *Iterator) Next() interface{} {
	if !i.HasNext() {
		panic("TimeCodeIterator completed")
	}
	ret := i.tc.TimeCode()
	i.tc.Next()
	return ret
}

func (tc *TimeCode) runUntil(tcf TimeCodeFragment) *Iterator {
	// Add 1 frame as end is the TimeCode of the frame after the iterator
	return &Iterator{
		tc:  tc,
		end: tcf.Add(0, 0, 0, 0, 1),
	}
}

func (tc *TimeCode) ForFrames(count int) *Iterator {
	return tc.runUntil(tc.TimeCode().AddFrames(count))
}

func (tc *TimeCode) Until(s string) (*Iterator, error) {
	tcf, err := ParseTimeCode(s, tc.FrameRate())
	if err != nil {
		return nil, err
	}

	// check for crossing midnight
	if tcf.Before(tc.TimeCode()) {
		tcf.day = tcf.day + 1
	}

	return tc.runUntil(tcf), nil
}
