package layout

import (
	"github.com/peter-mount/go-anim/layout"
	"image"
	"strconv"
)

type ContainerBuilder struct {
	ComponentBuilder
	comp layout.Container
}

func newContainerBuilder(parent any, builder *Builder, comp layout.Container) *ContainerBuilder {
	cb := &ContainerBuilder{
		ComponentBuilder: ComponentBuilder{
			comp:    comp,
			parent:  parent,
			builder: builder,
		},
		comp: comp,
	}
	cb.ComponentBuilder.this = cb
	return cb
}

func (b *ContainerBuilder) AddComponent(name string, comp layout.Component) (any, error) {
	if err := b.builder.add(name, comp); err != nil {
		return nil, err
	}
	b.comp.Add(comp)
	return newComponentBuilder(b, b.builder, comp), nil
}

func (b *ContainerBuilder) addContainer(comp layout.Container) (any, error) {
	b.builder.seq++
	comp.SetType(comp.GetType() + strconv.Itoa(b.builder.seq))
	b.comp.Add(comp)
	return newContainerBuilder(b, b.builder, comp), nil
}

func (b *ContainerBuilder) ColScaleContainer(scales ...float64) (any, error) {
	return b.addContainer(layout.ColScaleContainer(scales...))
}

func (b *ContainerBuilder) FixedContainer(width, height int) (any, error) {
	return b.addContainer(layout.FixedContainer(image.Rect(0, 0, width, height)))
}

func (b *ContainerBuilder) Image(name string) (any, error) {
	return b.AddComponent(name, layout.NewImage())
}

func (b *ContainerBuilder) RowContainer() (any, error) {
	return b.addContainer(layout.RowContainer())
}

func (b *ContainerBuilder) Text(name, format string, args ...any) (any, error) {
	return b.AddComponent(name, layout.NewText(format, args...))
}
