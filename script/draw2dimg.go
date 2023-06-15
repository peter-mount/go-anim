package script

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"image/draw"
)

type Draw2Dimg struct{}

func (_ Draw2Dimg) NewGraphicContext(img draw.Image) *draw2dimg.GraphicContext {
	return draw2dimg.NewGraphicContext(img)
}
