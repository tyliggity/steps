package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/public-steps/kubectl/base/events/get"
)

func run() (int, error) {
	eventGet, err := get.NewGetEvents(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(eventGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing get-events step: %v", err)
	}

	os.Exit(exitCode)
}
