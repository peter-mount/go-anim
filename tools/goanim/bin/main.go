package main

import (
	"fmt"
	"github.com/peter-mount/go-anim/tools/goanim"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(&goanim.Anim{}); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
