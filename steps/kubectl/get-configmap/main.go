package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/steps/kubectl/base"
	"github.com/stackpulse/steps/kubectl/base/configmaps/get"
)

func run() (int, error) {
	configmapGet, err := get.NewGetConfigmap(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(configmapGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing get-configmap step: %v", err)
	}

	os.Exit(exitCode)
}
