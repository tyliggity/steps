package main

import (
	"fmt"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/istio/base"
	jsonFormatter "github.com/stackpulse/steps-sdk-go/formats/json"
	"github.com/stackpulse/steps-sdk-go/parsers"
	"github.com/stackpulse/steps-sdk-go/step"
)

const proxyStatusArgument = "proxy-status"

type Args struct {
	base.Args
	Pod string `env:"POD" envDefault:""`
}

type ProxyStatus struct {
	args Args
	*base.IstioCtlBase
	kubeconfigPath string
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

	args = append(args, proxyStatusArgument)

	if p.args.Pod != "" {
		args = append(args, p.args.Pod)
	}

	exitCode, output, executionError := base.Execute(p.args.BinaryName, args)
	if executionError != nil {
		return exitCode, output, executionError
	}

	// when pod is empty it returns a table
	if p.args.Pod == "" {
		parsedOutput, err := parsers.ParseTableToJSON(string(output))
		if err == nil {
			return exitCode, parsedOutput, executionError
		}
	}

	encodedOutput, err := jsonFormatter.FormatAsJSONOutput(output)
	if err == nil {
		return exitCode, encodedOutput, executionError
	}

	// fallback

	return exitCode, output, executionError
}

func main() {
	step.Run(&ProxyStatus{})
}
