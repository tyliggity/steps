package base

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const DataRootJSONName = "data"

type Args struct {
	Host            string `env:"HOST,required"`
	Password        string `env:"PASSWORD,required"`
	Port            int    `env:"PORT,required"`
	User            string `env:"USER,required"`
	DB              string `env:"DB,required"`
	Pretty          bool   `env:"PRETTY" envDefault:"false"`
	FieldSeparator  string `env:"FIELD_SEPERATOR" envDefault:";;-;;"`
	RecordSeparator string `env:"RECORD_SEPERATOR" envDefault:";;=;;"`
}

type PsqlCommand struct {
	*Args
	sshArgs    *SSHArgs
	dockerArgs *DockerExecArgs
}

func NewPsqlCommand() (*PsqlCommand, error) {
	sshArgs, err := parseSSHArgs()
	if err != nil {
		return nil, fmt.Errorf("parse ssh args: %w", err)
	}

	dockerArgs, err := parseDockerExecArgs()
	if err != nil {
		return nil, fmt.Errorf("parse docker args: %w", err)
	}

	args := &Args{}
	if err := env.Parse(args); err != nil {
		return nil, err
	}

	if err := os.Setenv("PGPASSWORD", args.Password); err != nil {
		return nil, fmt.Errorf("can't set PGPASSWORD env: %w", err)
	}

	return &PsqlCommand{Args: args, sshArgs: sshArgs, dockerArgs: dockerArgs}, nil
}

func (p *PsqlCommand) RunPsqlCommand(extraArgs []string) ([]byte, int, error) {
	if outputBuffer, err := p.runSSHTunnel(); err != nil {
		return outputBuffer, 1, fmt.Errorf("running SSH tunnel: %w", err)
	}

	args := []string{"-A", // Don't align result
		"-F", p.FieldSeparator, // Seperation between fields in psql output
		"-R", p.RecordSeparator, // Seperation between records in psql output
		"-w",         // Don't prompt for password if missing
		"-h", p.Host, // Host
		"-p", strconv.Itoa(p.Port), // Port
		"-U", p.User, // User
		"-d", p.DB, // DB name
	}
	args = append(args, extraArgs...)

	var err error
	execName, args := p.getDockerExecCommand("psql", args)
	execName, args, err = p.getSSHExecCommand(execName, args)
	if err != nil {
		return nil, step.ExitCodeFailure, fmt.Errorf("get ssh exec command: %w", err)
	}

	cmd := exec.Command(execName, args...)

	log.Debugln("About to run %q with args: %#v", execName, args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return output, exiterr.ExitCode(), err
		}
		return output, 1, err
	}

	return output, 0, nil
}

func (p *PsqlCommand) ParseOutputJSON(output []byte) *gabs.Container {
	lines := strings.Split(string(output), "\n")
	startIndex := 0
	for i, line := range lines {
		if strings.Contains(line, p.FieldSeparator) {
			startIndex = i
			break
		}
	}
	lines = lines[startIndex:]
	output = []byte(strings.Join(lines, "\n"))

	ql := strings.Split(strings.TrimSpace(string(output)), p.RecordSeparator)
	columns := strings.Split(ql[0], p.FieldSeparator)
	for i, col := range columns {
		columns[i] = col
	}
	rootObj := gabs.New()

	// Starting from the second line (the actual rows) run until 1 line before the end of the output (1 line before end of output is the total rows)
	for i := 1; i < len(ql)-1; i++ {
		rowObj := gabs.New()
		line := strings.Split(ql[i], p.FieldSeparator)
		for j := 0; j < len(columns); j++ {
			if len(line) <= j {
				log.Logln("Failed parsing line #%d, can't find column index #%d within input (line: %s)", i, j, ql[i])
				continue
			}
			val := line[j]
			rowObj.Set(val, columns[j])
		}
		rootObj.ArrayAppend(rowObj, DataRootJSONName)
	}
	return rootObj
}

func (p *PsqlCommand) ParseOutput(output []byte) []byte {
	if !env.FormatterIs(env.JsonFormat) {
		return output
	}

	log.Debugln("Original output:\n%s\n", string(output))

	gc := p.ParseOutputJSON(output)

	if p.Pretty {
		return gc.BytesIndent("", "  ")
	}
	return gc.Bytes()
}
