package layout

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/peter-mount/go-anim/graph/resize"
	"image"
)

type Image struct {
	BaseComponent
	image image.Image
}

func NewImage() *Image {
	i := &Image{
		BaseComponent: BaseComponent{
			Type: "Image",
		},
	}
	i.painter = i.paint
	return i
}

func (i *Image) Layout(_ draw2d.GraphicContext) bool {
	//if !i.updateRequired {
	//	return false
	//}

	cb := i.Bounds()
	if i.image != nil && cb.Dx() == 0 && cb.Dy() == 0 {
		ib := i.image.Bounds()
		w := min(ib.Dx(), ib.Dy())
		cb.Max.X = cb.Min.X + w
		cb.Max.Y = cb.Min.Y + w
		i.SetBounds(cb)
	}

	i.updateRequired = false
	return true
}

func (i *Image) paint(gc *draw2dimg.GraphicContext) {
	if i.image == nil {
		return
	}

	cb := i.Bounds()

	// Resize the image to fit the component
	img := i.image
	ib := img.Bounds()

	if !cb.Eq(ib) {
		img = resize.Resize(uint(cb.Dx()), 0 /*uint(cb.Dy())*/, img, resize.NearestNeighbor)
	}

	gc.DrawImage(img)
}

func (i *Image) SetImage(img image.Image) {
	i.image = img
	i.updateRequired = true
}

func (i *Image) Image() image.Image {
	return i.image
}
