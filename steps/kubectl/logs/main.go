package main

import (
	"fmt"
	"os"

	"github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/public-steps/kubectl/base/logs/get"
)

func run() (int, error) {
	logsGet, err := get.NewGetLogs(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(logsGet)
}

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing logs step: %v", err)
	}

	os.Exit(exitCode)
}
