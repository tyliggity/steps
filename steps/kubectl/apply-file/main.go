package main

import (
	"fmt"
	"github.com/stackpulse/public-steps/kubectl/base"
	kctApply "github.com/stackpulse/public-steps/kubectl/base/apply"
	"os"
)

func run() (int, error) {
	applyKct, err := kctApply.NewApply(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(applyKct)
}

func main() {
	exitCode, err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing delete step: %v", err)
	}

	os.Exit(exitCode)
}
