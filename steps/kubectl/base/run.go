package base

import (
	"fmt"
	"os"
)

type Getter interface {
	Get() (output []byte, exitCode int, err error)
}

type Parsable interface {
	Parse(output []byte) (string, error)
}

type Deleter interface {
	Delete() (output []byte, exitCode int, err error)
}

type Runner interface {
	Run() (output []byte, exitCode int, err error)
}

func Run(cmd interface{}) (int, error) {
	var output []byte
	var exitCode int
	var err error

	switch t := cmd.(type) {
	case Getter:
		output, exitCode, err = t.Get()
	case Deleter:
		output, exitCode, err = t.Delete()
	case Runner:
		output, exitCode, err = t.Run()
	default:
		return 1, fmt.Errorf("unknown cmd type %T", t)
	}

	outputString := string(output)
	defer func(output *string) {
		os.Stdout.WriteString(*output)
	}(&outputString)

	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}
		return exitCode, err
	}

	parsable, ok := cmd.(Parsable)
	if !ok {
		return exitCode, nil
	}

	parsed, err := parsable.Parse(output)
	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}
		return exitCode, err
	}
	outputString = parsed
	return exitCode, nil
}
