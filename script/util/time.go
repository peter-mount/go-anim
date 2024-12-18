package util

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	RFC3339     = "2006-01-02T15:04:05Z07:00" // for convenience
	RFC3339Zulu = "2006-01-02T15:04:05Z"
	RFC3339NoTC = "2006-01-02T15:04:05"
	TIMESTAMP   = "20060102150405"
	TIMESTAMP2  = "060102150405"
)

var (
	timeFormats = []string{
		TIMESTAMP,
		TIMESTAMP2,
		RFC3339NoTC,
		RFC3339Zulu,
		RFC3339NoTC,
		RFC3339,
	}
)

func (u *Util) ParseTime(s string) time.Time {
	return u.ParseTimeIn(s, time.Local)
}

func (_ *Util) ParseTimeIn(s string, loc *time.Location) time.Time {

	// Parse time using one of our formats
	for _, tf := range timeFormats {
		if t, err := time.ParseInLocation(tf, s, loc); err == nil {
			return t
		}
	}

	// Unix time
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(i, 0)
	}

	return time.Time{}
}

func (u *Util) TimeFromFileName(s string) time.Time {
	return u.TimeFromFileNameIn(s, time.Local)
}

func (u *Util) TimeFromFileNameIn(s string, loc *time.Location) time.Time {
	_, s = filepath.Split(s)
	s = strings.TrimSuffix(s, filepath.Ext(s))
	return u.ParseTimeIn(s, loc)
}
