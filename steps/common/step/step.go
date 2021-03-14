package step

import (
	"fmt"
	"os"
)

const (
	ExitCodeOK      = 0
	ExitCodeFailure = 1
)

type Outputs struct {
	Object interface{} `json:"api_object,omitempty"`
}

// An interface describing a general purpose stackpulse step
type Step interface {
	Init() error
	Run() (exitCode int, output []byte, err error)
}

// A helper method that execute a given step according to step lifecycle
func Run(step Step) {
	err := step.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed init step arguments, %v\n", err)
		os.Exit(ExitCodeFailure)
	}

	exitCode, output, err := step.Run()
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("Step error: %s", err.Error())))
		if output != nil {
			os.Stderr.Write([]byte("\n"))
			os.Stderr.Write(output)
		}
		os.Exit(exitCode)
	}

	os.Stdout.Write(output)
}
