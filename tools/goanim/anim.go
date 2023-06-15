package goanim

import (
	_ "github.com/peter-mount/go-anim/script"
	_ "github.com/peter-mount/go-script/stdlib"
	_ "github.com/peter-mount/go-script/stdlib/io"
	_ "github.com/peter-mount/go-script/stdlib/math"
	_ "github.com/peter-mount/go-script/stdlib/time"
	"github.com/peter-mount/go-script/tools/goscript"
)

type Anim struct {
	_ *goscript.Script `kernel:"inject"`
}

func (a *Anim) Start() error {
	return nil
}
