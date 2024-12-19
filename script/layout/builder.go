package layout

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/peter-mount/go-anim/layout"
	"image"
)

type Builder struct {
	common
}
type common struct {
	components map[string]layout.Component
	root       layout.Container
}

func (*Package) New(width, height int) any {
	b := &Builder{
		common: common{
			components: make(map[string]layout.Component),
			root:       layout.FixedContainer(image.Rect(0, 0, width, height)),
		},
	}
	return newContainerBuilder(b, b, b.root)
}

func (b *Builder) add(n string, component layout.Component) error {
	if n != "" {
		if _, exists := b.components[n]; exists {
			return fmt.Errorf("component '%s' already exists", n)
		}
		b.components[n] = component
		component.SetType(n)
	}
	return nil
}

func (b *Builder) Build() *Layout {
	return &Layout{
		common: b.common,
	}
}

type Layout struct {
	common
}

func (l *Layout) Get(name string) layout.Component {
	return l.components[name]
}

func (l *Layout) Layout(ctx draw2d.GraphicContext) bool {
	r := l.root.Layout(ctx)
	if r {
		r = l.root.Layout(ctx)
	}
	return r
}

func (l *Layout) Draw(context draw2d.GraphicContext) {
	l.root.Draw(context)
}
