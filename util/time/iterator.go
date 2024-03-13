package time

import "github.com/peter-mount/go-kernel/v2/util"

// iterator is returned by TimeCode for advancing it for a set period.
// Each call to Next() will advance the TimeCode.
type iterator struct {
	tc  *TimeCode
	end TimeCodeFragment
}

func (i *iterator) HasNext() bool {
	return i.tc.TimeCode().Before(i.end)
}

func (i *iterator) Next() TimeCodeFragment {
	if !i.HasNext() {
		panic("TimeCodeIterator completed")
	}
	ret := i.tc.TimeCode()
	i.tc.Next()
	return ret
}

func (i *iterator) ForEach(f func(TimeCodeFragment)) {
	for i.HasNext() {
		f(i.Next())
	}
}

func (i *iterator) ForEachAsync(f func(TimeCodeFragment)) {
	i.ForEachAsync(f)
}

func (i *iterator) ForEachFailFast(f func(TimeCodeFragment) error) error {
	for i.HasNext() {
		if err := f(i.Next()); err != nil {
			return err
		}
	}
	return nil
}

func (i *iterator) Iterator() util.Iterator[TimeCodeFragment]        { return i }
func (i *iterator) ReverseIterator() util.Iterator[TimeCodeFragment] { return i }

func (tc *TimeCode) runUntil(tcf TimeCodeFragment) util.Iterator[TimeCodeFragment] {
	// Add 1 frame as end is the TimeCode of the frame after the iterator
	return &iterator{
		tc:  tc,
		end: tcf.Add(0, 0, 0, 0, 1),
	}
}

func (tc *TimeCode) ForFrames(count int) util.Iterator[TimeCodeFragment] {
	return tc.runUntil(tc.TimeCode().AddFrames(count))
}

func (tc *TimeCode) Until(s string) (util.Iterator[TimeCodeFragment], error) {
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
