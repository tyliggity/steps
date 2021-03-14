package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	gcloudAuthTmpFile = "/tmp/auth.json"
	kubeConfigTmpFile = "/tmp/kubeconfig"
)

func decodeAndWrite(base64Encoded, destination string) error {
	decoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, decoded, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// Running gcloud activate-service-account which is essential before running kubectl in gke custer with kubeconfig
func (fc *FluxCommand) runGcloudAuth() error {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return fmt.Errorf("gcloud command not found (are you running <step>-gke variation?)")
	}

	if err := decodeAndWrite(fc.args.GcloudAuth, gcloudAuthTmpFile); err != nil {
		return fmt.Errorf("can't write decoded gcloud auth: %w", err)
	}
	if fc.args.GcloudAuth == "" {
		if _, err := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", gcloudAuthTmpFile).CombinedOutput(); err != nil {
			return fmt.Errorf("can't run gcloud activate-service-account command: %w", err)
		}
	}
	return nil
}

func (fc *FluxCommand) getK8sAuth() (string, error) {
	if err := decodeAndWrite(fc.args.KubeConfigContent, kubeConfigTmpFile); err != nil {
		return "", fmt.Errorf("can't write decoded kubeconfig: %w", err)
	} else {
		return fmt.Sprintf(kubeConfigTmpFile), nil
	}
	// If kubeconfig content has specified, so it muse be local token
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", fmt.Errorf("can't read k8s local token file: %w", err)
	}
	return string(token), nil
}
