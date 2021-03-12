package exec

import (
	"github.com/stackpulse/public-steps/common/step"
	"os/exec"
)

func Execute(command string, args []string) (output []byte, exitCode int, err error) {
	cmd := exec.Command(command, args...)

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return cmdOutput, exiterr.ExitCode(), err
		}
		return cmdOutput, step.ExitCodeFailure, err
	}

	return cmdOutput, step.ExitCodeOK, nil
}
