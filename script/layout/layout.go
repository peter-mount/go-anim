package layout

import (
	"github.com/peter-mount/go-anim/layout"
	"github.com/peter-mount/go-script/packages"
	"image"
)

func init() {
	packages.RegisterPackage(&Package{})
}

type Package struct{}

func (*Package) RowContainer() layout.Container {
	return layout.RowContainer()
}

func (*Package) FixedContainer(rect image.Rectangle) layout.Container {
	return layout.FixedContainer(rect)
}

func (*Package) ColScaleContainer(scales ...float64) layout.Container {
	return layout.ColScaleContainer(scales...)
}

func (*Package) Image() layout.Component {
	return layout.NewImage()
}

func (*Package) Text(format string, args ...any) layout.Component {
	return layout.NewText(format, args...)
}
