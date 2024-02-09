package exr

import (
	"encoding/binary"
	"math"
)

type V2F struct {
	X, Y float32
}

func (v V2F) Bytes() []byte {
	b := binary.LittleEndian.AppendUint32(nil, math.Float32bits(v.X))
	b = binary.LittleEndian.AppendUint32(b, math.Float32bits(v.Y))
	return b
}
