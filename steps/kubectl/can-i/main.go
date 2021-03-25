package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/steps/kubectl/base"
	can_i "github.com/stackpulse/steps/kubectl/base/auth/can-i"
)

func run() (int, error) {
	cani, err := can_i.NewCanI(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(cani)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing can-i step: %v", err)
	}

	os.Exit(exitCode)
}
