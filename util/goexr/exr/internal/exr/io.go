package exr

import (
	"encoding/binary"
	"io"
)

var (
	order = binary.LittleEndian
)

func Read(in io.Reader, data any) error {
	return binary.Read(in, order, data)
}

func Write(w io.Writer, data any) error {
	return binary.Write(w, order, data)
}

func ReadNullTerminatedString[T ~string](in io.Reader, target *T) error {
	var buffer []byte
	for {
		var char byte
		if err := Read(in, &char); err != nil {
			return err
		}
		if char == 0x00 {
			break
		}
		buffer = append(buffer, char)
	}
	*target = T(buffer)
	return nil
}

func WriteNullTerminatedString[T ~string](w io.Writer, s T) error {
	b := AppendNullTerminatedString(nil, s)
	_, err := w.Write(b)
	return err
}

func AppendNullTerminatedString[T ~string](b []byte, s T) []byte {
	b = append(b, []byte(s)...)
	return append(b, 0x00)
}
