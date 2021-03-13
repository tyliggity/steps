package main

import (
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps/common/exec"
	"github.com/stackpulse/steps/common/log"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/influx/base"
)

type Args struct {
	base.Args
	Query string `env:"QUERY,required"`
}

type InfluxQuery struct {
	args Args
}

func (i *InfluxQuery) Init() error {
	err := envconf.Parse(&i.args)
	if err != nil {
		return err
	}

	return nil
}

func (i *InfluxQuery) buildCommand() []string {
	command := []string{"-database", i.args.Database}

	if i.args.SSL {
		command = append(command, "-ssl")
	}

	if i.args.UnsafeSSL {
		command = append(command, "-unsafeSsl")
	}

	if i.args.Username != "" {
		command = append(command, "-username", i.args.Username)
	}

	if i.args.Password != "" {
		command = append(command, "-password", i.args.Password)
	}

	command = append(command, "-port", i.args.Port, "-execute", i.args.Query, "-format", "json", "-pretty")

	log.Debug("influx cli command arguments", command)

	return command
}

func (i *InfluxQuery) Run() (int, []byte, error) {
	output, exitCode, err := exec.Execute(i.args.BinaryName, i.buildCommand())

	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}

		return exitCode, output, err
	}

	return exitCode, output, err
}

func main() {
	step.Run(&InfluxQuery{})
}
