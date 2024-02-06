package exr

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func ReadAttributeName(in io.Reader, target *AttributeName) error {
	return ReadNullTerminatedString(in, target)
}

const (
	AttributeNameChannels           AttributeName = "channels"
	AttributeNameCompression        AttributeName = "compression"
	AttributeNameDataWindow         AttributeName = "dataWindow"
	AttributeNameDisplayWindow      AttributeName = "displayWindow"
	AttributeNameLineOrder          AttributeName = "lineOrder"
	AttributeNamePixelAspectRatio   AttributeName = "pixelAspectRatio"
	AttributeNameScreenWindowCenter AttributeName = "screenWindowCenter"
	AttributeNameScreenWindowWidth  AttributeName = "screenWindowWidth"
)

type AttributeName string

func ReadAttributeType(in io.Reader, target *AttributeType) error {
	return ReadNullTerminatedString(in, target)
}

const (
	AttributeTypeChannelList AttributeType = "chlist"
	AttributeTypeCompression AttributeType = "compression"
	AttributeTypeBox2i       AttributeType = "box2i"
	AttributeTypeLineOrder   AttributeType = "lineOrder"
	AttributeTypeFloat       AttributeType = "float"
	AttributeTypeV2f         AttributeType = "v2f"
)

type AttributeType string

type BytesAttribute interface {
	Bytes() []byte
}

func WriteAttribute(w io.Writer, n AttributeName, t AttributeType, a any) error {
	if v, ok := a.(BytesAttribute); ok {
		b := v.Bytes()
		return WriteAttributeBytes(w, n, t, int32(len(b)), b)
	}

	if b, ok := a.([]byte); ok {
		return WriteAttributeBytes(w, n, t, int32(len(b)), b)
	}

	if i, ok := a.(*float32); ok {
		b := binary.LittleEndian.AppendUint32(nil, math.Float32bits(*i))
		return WriteAttributeBytes(w, n, t, int32(len(b)), b)
	}

	return fmt.Errorf("unsupported Attribute %q %q %T", n, t, a)
}

func WriteAttributeBytes(w io.Writer, n AttributeName, t AttributeType, s int32, b []byte) error {
	err := WriteAttributeHeader(w, n, t, s)
	if err == nil {
		_, err = w.Write(b)
	}
	return err
}

// WriteAttributeHeader writes an attribute minus its data. It's used for when the size is known
// but the data is not available as a single []byte
func WriteAttributeHeader(w io.Writer, n AttributeName, t AttributeType, s int32) error {
	err := WriteNullTerminatedString(w, n)
	if err == nil {
		err = WriteNullTerminatedString(w, t)
	}
	if err == nil {
		err = Write(w, &s)
	}
	return err
}
