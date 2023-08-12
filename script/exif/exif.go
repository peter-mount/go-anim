package exif

import (
	"github.com/peter-mount/go-script/packages"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"io"
	"math/big"
)

func init() {
	packages.Register("exif", &Exif{})
}

type Exif struct {
}

type Tags map[string]*tiff.Tag

func (t *Tags) Contains(n string) bool {
	e, exists := (*t)[n]
	// Note: e!=nil check to ensure we don't have a nil value
	return exists && e != nil
}

func (t *Tags) get(n string) (*tiff.Tag, bool) {
	e, exists := (*t)[n]
	return e, exists
}

// Count returns the number of values in the named tag.
// Returns 0 if the tag does not exist, otherwise a positive number
func (t *Tags) Count(n string) int {
	if s, exists := t.get(n); exists {
		return int(s.Count)
	}
	return 0
}

func (t *Tags) Rat(i int, n string, d1, d2 int64) (*big.Rat, error) {
	if s, exists := t.get(n); exists {
		return s.Rat(i)
	}
	return big.NewRat(d1, d2), nil
}

func (t *Tags) Rat2(i int, n string, d1, d2 int64) (int64, int64, error) {
	if s, exists := t.get(n); exists {
		return s.Rat2(i)
	}
	return d1, d2, nil
}

func (t *Tags) Int(i int, n string, d int64) (int64, error) {
	if s, exists := t.get(n); exists {
		if f, err := s.Int64(i); err == nil {
			return f, nil
		}
		if f, err := s.Float(i); err == nil {
			return int64(f), nil
		}
	}
	return d, nil
}

func (t *Tags) Float(i int, n string, d float64) (float64, error) {
	if s, exists := t.get(n); exists {
		if f, err := s.Float(i); err == nil {
			return f, nil
		}
		if f, err := s.Int64(i); err == nil {
			return float64(f), nil
		}
	}
	return d, nil
}

func (t *Tags) String(n, d string) string {
	if s, exists := t.get(n); exists {
		// If the tag is a String then return it directly
		if s.Format() == tiff.StringVal {
			str, _ := s.StringVal()
			return str
		}

		// Default to String which will convert int/float/rat to a string for us
		return s.String()
	}
	return d
}

type exifTagsBuilder Tags

func (m *exifTagsBuilder) Walk(name exif.FieldName, tag *tiff.Tag) error {
	(*m)[string(name)] = tag
	return nil
}

func (_ Exif) Decode(r io.Reader) (*Tags, error) {
	b := make(exifTagsBuilder)

	x, err := exif.Decode(r)
	if err == nil {
		err = x.Walk(&b)
	}
	return (*Tags)(&b), err
}
