package goanim

import (
	"github.com/llgcode/draw2d"
	_ "github.com/peter-mount/go-anim/script"
	_ "github.com/peter-mount/go-script/stdlib"
	_ "github.com/peter-mount/go-script/stdlib/fmt"
	_ "github.com/peter-mount/go-script/stdlib/io"
	_ "github.com/peter-mount/go-script/stdlib/math"
	_ "github.com/peter-mount/go-script/stdlib/time"
	"github.com/peter-mount/go-script/tools/goscript"
	"os"
	"path"
	"path/filepath"
)

type Anim struct {
	_ *goscript.Script `kernel:"inject"`
}

func (a *Anim) Start() error {

	// Set location of our fonts
	draw2d.SetFontFolder(path.Join(filepath.Dir(os.Args[0]), "../lib/font"))

	return nil
}
