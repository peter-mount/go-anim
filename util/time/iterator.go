package time

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
	i.tc.Next()
	return i.tc.TimeCode()
}

func (tc *TimeCode) runUntil(tcf TimeCodeFragment) *Iterator {
	return &Iterator{
		tc:  tc,
		end: tcf,
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
