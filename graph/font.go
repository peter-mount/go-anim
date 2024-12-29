package graph

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/util/font"
)

func SetFont(gc *draw2dimg.GraphicContext, s string) error {
	f, err := font.ParseFont(s)
	if err != nil {
		return err
	}

	gc.SetFontData(f.FontData())
	gc.SetFontSize(f.Size())
	return nil
}
