package main

import (
	"fmt"
	"github.com/peter-mount/go-anim/tools/build"
	_ "github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
		&build.Install{},
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
