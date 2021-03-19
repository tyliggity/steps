package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/public-steps/kubectl/base/nodes/get"
)

func run() (int, error) {
	nodesGet, err := get.NewGetNodes(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(nodesGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing get-nodes step: %v", err)
	}

	os.Exit(exitCode)
}
