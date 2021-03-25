package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/steps/kubectl/base"
	"github.com/stackpulse/steps/kubectl/base/version/get"
)

func run() (int, error) {
	versionGet, err := get.NewGetVersion(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(versionGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing version step: %v", err)
	}

	os.Exit(exitCode)
}
