package exr

import (
	"errors"
	"io"
)

const (
	Extension = "exr"
)

var (
	MagicSequence = [4]byte{0x76, 0x2F, 0x31, 0x01}
	magicSlice    = []byte{0x76, 0x2F, 0x31, 0x01}
)

func ReadMagic(in io.Reader, target *Magic) error {
	return Read(in, target)
}

type Magic [4]byte

func (m Magic) IsCorrect() bool {
	return m == MagicSequence
}

func WriteMagic(w io.Writer) error {
	n, err := w.Write(magicSlice)
	if err != nil {
		return err
	}
	if n != len(magicSlice) {
		return errors.New("wrote too short Magic")
	}
	return nil
}
