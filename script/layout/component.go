package layout

import (
	"github.com/peter-mount/go-anim/layout"
	"image/color"
)

type ComponentBuilder struct {
	comp    layout.Component
	parent  any
	this    any
	builder *Builder
}

func newComponentBuilder(parent any, builder *Builder, comp layout.Component) *ComponentBuilder {
	cb := &ComponentBuilder{
		comp:    comp,
		parent:  parent,
		builder: builder,
	}
	cb.this = cb
	return cb
}

func (b *ComponentBuilder) Align(s string) any {
	b.comp.Align(s)
	return b.this
}

func (b *ComponentBuilder) Fill(c color.Color) any {
	b.comp.Fill(c)
	return b.this
}

func (b *ComponentBuilder) Font(font string) any {
	b.comp.Font(font)
	return b.this
}

func (b *ComponentBuilder) LineWidth(w float64) any {
	b.comp.LineWidth(w)
	return b.this
}

func (b *ComponentBuilder) Stroke(c color.Color) any {
	b.comp.Stroke(c)
	return b.this
}

func (b *ComponentBuilder) StrokeFill(c color.Color) any {
	b.comp.StrokeFill(c)
	return b.this
}

func (b *ComponentBuilder) End() any {
	return b.parent
}
