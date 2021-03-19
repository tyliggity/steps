package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/public-steps/kubectl/base/pods/get"
)

func run() (int, error) {
	podsGet, err := get.NewGetPods(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(podsGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing get-pods step: %v", err)
	}

	os.Exit(exitCode)
}
