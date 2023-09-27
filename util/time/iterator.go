package time

type Iterator struct {
	tc  *TimeCode
	end TimeCodeFragment
}

func (i *Iterator) HasNext() bool {
	return i.tc.TimeCode().NotAfter(i.end)
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
