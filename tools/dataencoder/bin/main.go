package main

import (
	"fmt"
	"github.com/peter-mount/go-anim/tools/dataencoder"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	err := kernel.Launch(
		&dataencoder.Build{},
		&dataencoder.Include{},
		&dataencoder.Lib{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
