package main

import (
	"os"

	"github.com/johnharris85/helm2csar/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
