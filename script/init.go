package script

import (
	"github.com/peter-mount/go-anim/script/io"
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

	packages.Register("ffmpeg", &io.FFMPeg{})
	packages.Register("jpeg", &io.JPEG{})
	packages.Register("png", &io.PNG{})
}
