package base

import (
	"fmt"
	"github.com/stackpulse/public-steps/common/env"
	"os"
	"strconv"
)

const DockerExecEnvKey = "USE_DOCKER_EXEC"

type DockerExecArgs struct {
	ContainerName           string `env:"DOCKER_CONTAINER_NAME,required"`
	UseSudo                 bool   `env:"DOCKER_USE_SUDO", envDefault:"false"`
	ContainerPostgresServer string `env:"DOCKER_CONTAINER_POSTGRES_SERVER" envDefault:"127.0.0.1"`
}

func parseDockerExecArgs() (*DockerExecArgs, error) {
	userDockerExecEnv := os.Getenv(DockerExecEnvKey)
	if userDockerExecEnv == "" {
		return nil, nil
	}

	useDockerExec, err := strconv.ParseBool(userDockerExecEnv)
	if err != nil {
		return nil, fmt.Errorf("parsing %s env as bool: %w", DockerExecEnvKey, err)
	}
	if !useDockerExec {
		return nil, nil
	}

	dockerExecArgs := &DockerExecArgs{}
	if err := env.Parse(dockerExecArgs); err != nil {
		return nil, fmt.Errorf("parsed docker exec args: %w", err)
	}

	// Overriding host to localhost (or other provided one) when using docker
	if err := os.Setenv(HostEnv, dockerExecArgs.ContainerPostgresServer); err != nil {
		return nil, fmt.Errorf("set %s env: %w", HostEnv, err)
	}

	return dockerExecArgs, nil
}

func (p *PsqlCommand) getDockerExecCommand(execName string, args []string) (string, []string) {
	if p.dockerArgs == nil {
		return execName, args
	}
	dockerArgs := []string{"exec", p.dockerArgs.ContainerName, execName}
	dockerArgs = append(dockerArgs, args...)

	dockerExecName := "docker"
	if p.dockerArgs.UseSudo {
		dockerExecName = "sudo"
		sudoArgs := []string{"-n", "docker"} // -n for not sudo not prompt for password if required
		dockerArgs = append(sudoArgs, dockerArgs...)
	}
	return dockerExecName, dockerArgs
}
