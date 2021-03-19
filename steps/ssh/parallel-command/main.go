package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"io/ioutil"
	"sync"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/exec"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	Username              string   `env:"USERNAME,required" envDefault:""`
	HostNames             []string `env:"HOSTNAMES,required" envDefault:""`
	Command               string   `env:"COMMAND,required" envDefault:""`
	PrivateKey            string   `env:"PRIVATE_KEY" envDefault:""`
	AWSSecretKey          string   `env:"AWS_SECRET_KEY" envDefault:""`
	AWSRegion             string   `env:"AWS_REGION" envDefault:""`
	StrictHostKeyChecking string   `env:"STRICT_HOST_KEY_CHECKING" envDefault:"no"`
	LogLevel              string   `env:"LOG_LEVEL" envDefault:"ERROR"`
}

type SSHResponse struct {
	Hostname string
	Success  bool
	Output   string
	Error    string
}

const (
	PrivateKeyPath = "/key.pem"
	MaxHosts       = 1000
)

type ParallelSSHCommand struct {
	args Args
}

func (p *ParallelSSHCommand) fetchAwsSecret(secretKey, region string) (string, error) {
	secretsService := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	if secretsService == nil {
		return "", fmt.Errorf("initialize AWS secrets manager")
	}

	result, err := secretsService.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     &secretKey,
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil || result.SecretString == nil {
		return "", fmt.Errorf("fetch aws secret: %w", err)
	}

	return *result.SecretString, nil
}

func (p *ParallelSSHCommand) Init() error {
	err := env.Parse(&p.args)
	if err != nil {
		return err
	}

	if len(p.args.HostNames) > MaxHosts {
		return fmt.Errorf("HOSTNAMES parameter exceeded the max hosts allowed which is %d", MaxHosts)
	}

	if p.args.PrivateKey != "" {
		err = ioutil.WriteFile(PrivateKeyPath, []byte(p.args.PrivateKey), 0644)
		if err != nil {
			return fmt.Errorf("write private key file: %w", err)
		}

		return nil
	}

	if p.args.AWSSecretKey != "" {
		privateKey, err := p.fetchAwsSecret(p.args.AWSSecretKey, p.args.AWSRegion)
		if err != nil {
			return fmt.Errorf("fetch aws secret: %w", err)
		}

		err = ioutil.WriteFile(PrivateKeyPath, []byte(privateKey), 0644)
		if err != nil {
			return fmt.Errorf("write AWS private key file: %w", err)
		}

		return nil
	}

	return fmt.Errorf("ssh private key is required")
}

type HostsOutputs struct {
	SSHResponses []SSHResponse
}

func (p *ParallelSSHCommand) Run() (int, []byte, error) {
	// Restrict key.pem file capabilities (for ssh usage)
	_, exitCode, err := exec.Execute("chmod", []string{"600", PrivateKeyPath})
	if err != nil {
		return exitCode, nil, fmt.Errorf("execute linux command chmod: err code: %d, error: %w", exitCode, err)
	}

	results := make(chan SSHResponse, len(p.args.HostNames))
	var wg sync.WaitGroup
	for i := 0; i < len(p.args.HostNames); i++ {
		wg.Add(1)
		go p.ExecuteSSH(&wg, p.args.HostNames[i], results)
	}

	wg.Wait()
	close(results)

	var hostsOutputs HostsOutputs
	for res := range results {
		hostsOutputs.SSHResponses = append(hostsOutputs.SSHResponses, res)
	}

	jsonResponses, err := json.Marshal(hostsOutputs)
	if err != nil {
		return exitCode, nil, fmt.Errorf("marshal responses failed, error: %w", err)
	}

	return step.ExitCodeOK, jsonResponses, err
}

func (p *ParallelSSHCommand) ExecuteSSH(wg *sync.WaitGroup, hostname string, output chan<- SSHResponse) {
	defer wg.Done()

	sshArgs := p.buildArgs(p.args.Username, hostname, p.args.Command, p.args.StrictHostKeyChecking, p.args.LogLevel)
	o, _, err := exec.Execute("ssh", sshArgs)
	if err != nil {
		output <- SSHResponse{
			Hostname: hostname,
			Success:  false,
			Error:    err.Error(),
		}
		return
	}

	output <- SSHResponse{
		Hostname: hostname,
		Success:  true,
		Output:   string(o),
	}
}

func (p *ParallelSSHCommand) buildArgs(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string) []string {
	args := []string{"-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking)}
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, "-i", PrivateKeyPath)
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, linuxCmd)

	return args
}

func main() {
	step.Run(&ParallelSSHCommand{})
}
