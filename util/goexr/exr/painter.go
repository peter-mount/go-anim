package exr

import (
	"github.com/golang/freetype/raster"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
)

type RGBAImagePainter struct {
	// Image is the image to compose onto.
	Image *RGBAImage
	Rect  image.Rectangle
	// Op is the Porter-Duff composition operator.
	Op draw.Op
	// cr, cg, cb and ca are the 16-bit color to paint the spans.
	cr, cg, cb, ca float32
}

// Paint satisfies the Painter interface.
func (r *RGBAImagePainter) Paint(ss []raster.Span, done bool) {
	b := r.Image.Bounds()
	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}

		// This code mimics drawGlyphOver in $GOROOT/src/image/draw/draw.go.
		const m = float32(0xffff)
		ma := float32(s.Alpha) / m
		a := (ma + r.ca) - (ma * r.ca)
		cr, cg, cb := r.cr*a, r.cg*a, r.cb*a
		if r.Op == draw.Over {
			for i := s.X0; i < s.X1; i++ {
				dr, dg, db, da := r.Image.GetRGBA(i, s.Y)
				a1 := da * (1 - a)
				r.Image.SetRGBA(i, s.Y,
					cr+(dr*a1),
					cg+(dg*a1),
					cb+(db*a1),
					a+da-(a*da),
				)
			}
		} else {
			for i := s.X0; i < s.X1; i++ {
				r.Image.SetRGBA(i, s.Y, cr, cg, cb, a)
			}
		}
	}
}

// SetColor sets the color to paint the spans.
func (r *RGBAImagePainter) SetColor(c color.Color) {
	const m = float32(0xffff)
	cr, cg, cb, ca := c.RGBA()
	r.cr = float32(cr) / m
	r.cg = float32(cg) / m
	r.cb = float32(cb) / m
	r.ca = float32(ca) / m
}

func NewRGBAImagePainter(m *RGBAImage) draw2dimg.Painter {
	return &RGBAImagePainter{Image: m}
}
