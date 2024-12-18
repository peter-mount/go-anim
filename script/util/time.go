package util

import (
	time2 "github.com/peter-mount/go-anim/util/time"
	"time"
)

func (_ *Util) ParseTime(s string) time.Time {
	return time2.ParseTimeIn(s, time.Local)
}

func (_ *Util) ParseTimeIn(s string, loc *time.Location) time.Time {
	return time2.ParseTimeIn(s, loc)
}

func (_ *Util) TimeFromFileName(s string) time.Time {
	return time2.TimeFromFileNameIn(s, time.Local)
}

func (_ *Util) TimeFromFileNameIn(s string, loc *time.Location) time.Time {
	return time2.TimeFromFileNameIn(s, loc)
}
