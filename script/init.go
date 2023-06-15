package script

import "github.com/peter-mount/go-script/packages"

func init() {
	packages.Register("colour", &Colour{})
	packages.Register("animGraphic", &Graph{})
	packages.Register("animUtil", &AnimUtil{})
}
