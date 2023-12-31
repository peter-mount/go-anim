package script

import (
	_ "github.com/peter-mount/go-anim/script/exif"
	_ "github.com/peter-mount/go-anim/script/io"
	"github.com/peter-mount/go-script/packages"
)

func init() {
	packages.Register("animGraphic", &Graph{})
	packages.Register("animUtil", &AnimUtil{})
	packages.Register("colour", &Colour{})
	packages.Register("draw2dimg", &Draw2Dimg{})
	packages.Register("image", newImage())

	packages.Register("graphFilter", Filter{})
	packages.Register("graphMapper", Mapper{})

}
