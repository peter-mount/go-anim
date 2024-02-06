package exr

import (
	"fmt"
	"io"
)

func ReadChannelList(in io.Reader, target *ChannelList) error {
	var channels []Channel
	for {
		var channel Channel
		if err := ReadNullTerminatedString(in, &channel.Name); err != nil {
			return fmt.Errorf("error reading channel name: %w", err)
		}
		if channel.Name == "" {
			break
		}
		if err := Read(in, &channel.PixelType); err != nil {
			return fmt.Errorf("error reading channel pixel type: %w", err)
		}
		if err := Read(in, &channel.Linear); err != nil {
			return fmt.Errorf("error reading channel linearity: %w", err)
		}
		var reserved [3]int8
		if err := Read(in, &reserved); err != nil {
			return fmt.Errorf("error reading channel reserved data: %w", err)
		}
		if err := Read(in, &channel.XSampling); err != nil {
			return fmt.Errorf("error reading channel x sampling: %w", err)
		}
		if err := Read(in, &channel.YSampling); err != nil {
			return fmt.Errorf("error reading channel y sampling: %w", err)
		}
		channels = append(channels, channel)
	}
	*target = ChannelList(channels)
	return nil
}

type ChannelList []Channel

func (l ChannelList) Bytes() []byte {
	var b []byte
	for _, c := range l {
		b = c.append(b)
	}
	return append(b, 0x00)
}

type Channel struct {
	Name      string
	PixelType PixelType
	Linear    bool
	XSampling int32
	YSampling int32
}

func (c Channel) append(b []byte) []byte {
	b = AppendNullTerminatedString(b, c.Name)
	b = order.AppendUint32(b, uint32(c.PixelType))
	if c.Linear {
		b = append(b, 1)
	} else {
		b = append(b, 0)
	}

	// 3 reserved bytes
	b = append(b, 0, 0, 0)

	b = order.AppendUint32(b, uint32(c.XSampling))
	b = order.AppendUint32(b, uint32(c.YSampling))
	return b
}

const (
	PixelTypeUint  PixelType = 0
	PixelTypeHalf  PixelType = 1
	PixelTypeFloat PixelType = 2
)

type PixelType int32

func (t PixelType) String() string {
	switch t {
	case PixelTypeUint:
		return "UINT"
	case PixelTypeHalf:
		return "HALF"
	case PixelTypeFloat:
		return "FLOAT"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}
