package main

import (
	"fmt"

	envconf "github.com/caarlos0/env/v6"
	jsonFormatter "github.com/stackpulse/steps-sdk-go/formats/json"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/istio/base"
)

const proxyConfigArgument = "proxy-config"

type Args struct {
	base.Args
	Command string `env:"COMMAND,required"`
	Pod     string `env:"POD,required"`
}

type ProxyStatus struct {
	args Args
	*base.IstioCtlBase
	kubeconfigPath string
}

var validCommands = []string{"bootstrap", "cluster", "endpoint", "lisener", "log", "route", "secret"}

func isValidCommand(c string) bool {
	for _, command := range validCommands {
		if c == command {
			return true
		}
	}

	return false
}

func (p *ProxyStatus) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	err = p.args.Validate()
	if err != nil {
		return err
	}

	if !isValidCommand(p.args.Command) {
		return fmt.Errorf("%s is not one of the valid commands", p.args.Command)
	}

	if !p.args.UseLocalToken {
		p.kubeconfigPath, err = base.DecodeKubeConfigContent(p.args.KubeconfigContent)
		if err != nil {
			return fmt.Errorf("failed decoding kubeconfig: %w", err)
		}
	}

	return nil
}

func (p *ProxyStatus) Run() (int, []byte, error) {
	args := p.args.BuildCommand(p.kubeconfigPath)

	args = append(args, proxyConfigArgument, p.args.Command, p.args.Pod, "-o", "json")

	exitCode, output, executionError := base.Execute(p.args.BinaryName, args)
	if exitCode == 0 {
		return exitCode, output, executionError
	}

	// in case there is an error the output is returned not as json but as error line string
	encodedOutput, err := jsonFormatter.FormatAsJSONOutput(output)
	if err == nil {
		return exitCode, encodedOutput, executionError
	}

	return exitCode, output, executionError
}

func main() {
	step.Run(&ProxyStatus{})
}
