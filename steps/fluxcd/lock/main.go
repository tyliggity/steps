package main

import (
	"encoding/json"
	"fmt"

	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Args struct {
	Name              string `env:"NAME,required"`
	Namespace         string `env:"NAMESPACE,required" envDefault:"default"`
	FluxURL           string `env:"FLUX_URL,required"`
	User              string `env:"USER,required" envDefault:""`
	KubeConfigContent string `env:"KUBECONFIG_CONTENT"`
	GcloudAuth        string `env:"GCLOUD_AUTH_CODE_B64"`
	Lock              bool   `env:"LOCK,required" envDefault:"False"`
}

type Output struct {
	User      string
	Operation string
	Namespace string
	Name      string
}

type FluxCommand struct {
	args Args
}

func (fc *FluxCommand) Init() error {
	return env.Parse(&fc.args)
}

func getKubernetesConfig(fc *FluxCommand) (*rest.Config, error) {
	var config *rest.Config
	var err error
	var authMethod string

	config, err = rest.InClusterConfig()
	if config != nil {
		return config, nil
	}
	if fc.args.GcloudAuth != "" {
		if err = fc.runGcloudAuth(); err != nil {
			return nil, fmt.Errorf("get gcloud auth method: %w", err)
		}
	}
	if fc.args.KubeConfigContent != "" {
		if authMethod, err = fc.getK8sAuth(); err != nil {
			return nil, fmt.Errorf("get k8sauth method: %w", err)
		}
	}
	if config, err = clientcmd.BuildConfigFromFlags("", authMethod); err != nil {
		return nil, fmt.Errorf("build config from flags: %w", err)
	}
	return config, nil
}

func (fc *FluxCommand) Run() (int, []byte, error) {
	var operation string
	var err error
	var ret []byte
	var config *rest.Config

	if config, err = getKubernetesConfig(fc); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get kubernetes config: %w", err)
	}
	k, err := kubernetes.NewForConfig(config)
	w := New(fc.args.FluxURL, k)
	if fc.args.Lock == true {
		err = w.Lock(fc)
		operation = "lock"

	} else {
		err = w.Unlock(fc)
		operation = "unlock"
	}
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("%s: %w", operation, err)
	}
	output := Output{User: fc.args.User, Operation: operation, Namespace: fc.args.Namespace, Name: fc.args.Name}
	if ret, err = json.Marshal(output); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal: %w", err)
	}
	return step.ExitCodeOK, ret, err
}

func main() {
	step.Run(&FluxCommand{})
}
