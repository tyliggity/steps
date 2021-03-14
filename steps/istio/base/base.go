package base

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stackpulse/public-steps/common/exec"
	"github.com/stackpulse/public-steps/common/log"
)

type Args struct {
	// Kubernetes configuration file
	KubeconfigContent string `env:"KUBECONFIG_CONTENT" envDefault:""`
	UseLocalToken     bool   `env:"USE_LOCAL_TOKEN" envDefault:"true"`

	// Istio system namespace
	IstioNamespace string `env:"ISTIO_NAMESPACE" envDefault:"istio-system"`
	// The name of the kubeconfig context to use
	Context string `env:"CONTEXT"`
	// config namespace
	Namespace  string `env:"NAMESPACE"`
	BinaryName string `env:"BINARY_NAME" envDefault:"istioctl"`
}

const KubeconfigFile = "/tmp/.kube-config"

type IstioCtlBase struct {
	args Args
}

// Decode kubeconfig base64 to a file and return the file path
func DecodeKubeConfigContent(content string) (string, error) {
	kubeconfig, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("failed decoding kubeconfig: %w", err)
	}

	if err := ioutil.WriteFile(KubeconfigFile, kubeconfig, os.FileMode(0644)); err != nil {
		return "", fmt.Errorf("failed creating kubeconfig file: %w", err)
	}

	return KubeconfigFile, nil
}

func (a *Args) Validate() error {
	if !a.UseLocalToken && a.KubeconfigContent == "" {
		return fmt.Errorf("neither USE_LOCAL_TOKEN nor KUBECONFIG_TOKEN was specified")
	}

	if a.KubeconfigContent != "" && a.UseLocalToken {
		return fmt.Errorf("USE_LOCAL_TOKEN and KUBECONFIG_CONTENT can't be both specified")
	}

	return nil
}

func (a *Args) BuildCommand(kubeconfigFile string) []string {
	command := []string{}

	if !a.UseLocalToken {
		command = append(command, "--kubeconfig", kubeconfigFile)
	}

	if a.IstioNamespace != "" {
		command = append(command, "--istioNamespace", a.IstioNamespace)
	}

	if a.Namespace != "" {
		command = append(command, "--namespace", a.Namespace)
	}

	if a.Context != "" {
		command = append(command, "--context", a.Context)
	}

	log.Debug("istioctl cli command arguments", command)

	return command
}

func Execute(binary string, args []string) (int, []byte, error) {
	output, exitCode, err := exec.Execute(binary, args)

	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}

		return exitCode, output, err
	}

	return exitCode, output, err
}
