package goanim

import (
	_ "github.com/peter-mount/go-script/stdlib"
	_ "github.com/peter-mount/go-script/stdlib/math"
	"github.com/peter-mount/go-script/tools/goscript"
)

type Anim struct {
	_ *goscript.Script `kernel:"inject"`
}

func (a *Anim) Start() error {
	return nil
}
