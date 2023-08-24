package clip

import (
	"encoding/binary"
	"github.com/peter-mount/go-anim/util/time"
	"io"
)

// Write saves a clip to a Writer.
//
// Important: The format of this is not guaranteed to remain the same
// between versions.
// Persistence of Clip's is only for temporary persistence only, e.g.
// Parsing a timeline as part of the workflow.
func (c *Clip) Write(w io.Writer) error {
	err := c.timeCode.Write(w)
	if err == nil {
		err = binary.Write(w, binary.BigEndian, len(c.frames))
	}
	if err != nil {
		return err
	}

	for _, f := range c.frames {
		if err := f.write(w); err != nil {
			return err
		}
	}
	return nil
}

func ReadClip(r io.Reader) (c *Clip, err error) {
	c = &Clip{}

	c.timeCode, err = time.ReadTimeCode(r)
	if err == nil {
		var l int
		err = binary.Read(r, binary.BigEndian, &l)

		for i := 0; i < l && err == nil; i++ {
			var f Frame
			err = readFrame(r, &f)
			c.frames = append(c.frames, f)
		}
	}

	return c, err
}

func (f Frame) write(w io.Writer) error {
	err := f.timeCode.Write(w)
	if err == nil {
		err = binary.Write(w, binary.BigEndian, f.name)
	}
	return err
}

func readFrame(r io.Reader, f *Frame) error {
	tc, err := time.ReadTimeCodeFragment(r)
	if err != nil {
		return err
	}
	f.timeCode = tc
	return binary.Read(r, binary.BigEndian, &f.name)
}
