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

func (i *Image) Layout(ctx draw2d.GraphicContext) bool {
	//if !i.updateRequired {
	//	return false
	//}

	i.updateRequired = false

	cb := i.Bounds()
	if i.image != nil {
		if cb.Dx() == 0 && cb.Dy() == 0 {
			ib := i.image.Bounds()
			w := min(ib.Dx(), ib.Dy())
			cb.Max.X = cb.Min.X + w
			cb.Max.Y = cb.Min.Y + w
			i.SetBounds(cb)
		}

		img := i.image
		ib := img.Bounds()

		// Now use the inset bounds for the image
		cb = i.InsetBounds()

		if !cb.Eq(ib) {
			// FIXME work out new size without creating an image then throwing it away
			img = resize.Resize(uint(cb.Dx()), 0, img, resize.NearestNeighbor)
			ib = img.Bounds()
			cb.Max = cb.Min.Add(image.Pt(ib.Dx(), ib.Dy()))
			i.SetBounds(cb)
		}
	}

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
		img = resize.Resize(uint(cb.Dx()-i.insetX-i.insetX), 0, img, resize.NearestNeighbor)
		ib = img.Bounds()
		cb.Max = cb.Min.Add(image.Pt(ib.Dx()+(i.insetX<<1), ib.Dy()+i.insetMinY+i.insetMaxY))
		i.SetBounds(cb)
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
