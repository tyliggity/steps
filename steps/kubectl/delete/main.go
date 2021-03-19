package main

import (
	"fmt"
	"os"

	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/public-steps/kubectl/base"
	kctDelete "github.com/stackpulse/public-steps/kubectl/base/delete"
	"github.com/stackpulse/steps-sdk-go/env"
)

func run() (int, error) {
	deleteKct, err := kctDelete.NewDelete(nil)
	if err != nil {
		return 1, err
	}

	return base.Run(deleteKct)
}

func main() {
	exitCode, err := run()

	gc := gabs.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing delete step: %v", err)
		gc.Set(false, "success")
	} else {
		gc.Set(true, "success")
	}

	if env.FormatterIs(env.JsonFormat) {
		fmt.Printf("\n%s\n", gc.String())
	}

	os.Exit(exitCode)
}
