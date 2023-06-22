package main

import (
	"fmt"
	"github.com/peter-mount/go-anim/tools/build"
	_ "github.com/peter-mount/go-build/tools/build"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	err := kernel.Launch(
		&build.Include{},
		&build.Lib{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
