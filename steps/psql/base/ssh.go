package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
)

// For overriding those env in SSH tunnel mode, pay attention to change those as well if you change the env struct tags
const HostEnv = "HOST"
const PortEnv = "PORT"

const SleepBetweenCheckingTunnel = 50 * time.Millisecond
const SSHTunnelListenAdder = "127.0.0.1"
const SSHTunnelListenPort = "5432"
const SSHEnvKey = "USE_SSH"

type SSHArgs struct {
	SSHHost              string        `env:"SSH_HOST,required"`
	SSHKey               string        `env:"SSH_KEY,required"`
	SSHUser              string        `env:"SSH_USER,required"`
	SSHPort              int           `env:"SSH_PORT" envDefault:"22"`
	TunnelMode           bool          `env:"SSH_TUNNEL_MODE" envDefault:"false"`
	RemotePostgresServer string        `env:"SSH_REMOTE_POSTGRES_SERVER" envDefault:"127.0.0.1"`
	RemotePostgresPort   int           `env:"SSH_REMOTE_POSTGRES_PORT" envDefault:"5432"`
	WaitForTunnelTimeout time.Duration `env:"SSH_WAIT_FOR_TUNNEL_TIMEOUT" envDefault:"30s"`
}

func parseSSHArgs() (*SSHArgs, error) {
	useSShEnv := os.Getenv(SSHEnvKey)
	if useSShEnv == "" {
		return nil, nil
	}

	useSSh, err := strconv.ParseBool(useSShEnv)
	if err != nil {
		return nil, fmt.Errorf("parsing %s env as bool: %w", SSHEnvKey, err)
	}
	if !useSSh {
		return nil, nil
	}

	sshArgs := &SSHArgs{}
	if err := env.Parse(sshArgs); err != nil {
		return nil, fmt.Errorf("parse ssh args: %w", err)
	}

	if sshArgs.TunnelMode {
		// Overriding host and port to the tunnel host and port
		if err := os.Setenv(HostEnv, SSHTunnelListenAdder); err != nil {
			return nil, fmt.Errorf("set %s env: %w", HostEnv, err)
		}
		if err := os.Setenv(PortEnv, SSHTunnelListenPort); err != nil {
			return nil, fmt.Errorf("set %s env: %w", PortEnv, err)
		}
	}

	return sshArgs, nil
}

func (p *PsqlCommand) writeSSHKEyFile() (string, error) {
	tempFile, err := ioutil.TempFile("", "key*.pem")
	if err != nil {
		return "", fmt.Errorf("create tempfile for ssh key: %w", err)
	}
	if _, err := tempFile.Write([]byte(p.sshArgs.SSHKey)); err != nil {
		return "", fmt.Errorf("write key to temp file: %w", err)
	}
	return tempFile.Name(), nil
}

func (p *PsqlCommand) buildBaseSSHArgs(keyFile string, additionalParameters []string) []string {
	parameters := []string{
		"-i", keyFile, // Private key for connecting
		"-oStrictHostKeyChecking=no", // Don't ask to verify the host
		fmt.Sprintf("-oConnectTimeout=%d", int(p.sshArgs.WaitForTunnelTimeout.Seconds())), // Connection timeout
		"-4", // Connect on IPv4
		"-q",
		"-p", fmt.Sprintf("%d", p.sshArgs.SSHPort), // SSH Port
	}
	parameters = append(parameters, additionalParameters...)
	return append(parameters, fmt.Sprintf("%s@%s", p.sshArgs.SSHUser, p.sshArgs.SSHHost))
}

func (p *PsqlCommand) runSSHTunnel() ([]byte, error) {
	var b bytes.Buffer
	if p.sshArgs == nil || !p.sshArgs.TunnelMode {
		return []byte{}, nil
	}

	keyFile, err := p.writeSSHKEyFile()
	if err != nil {
		return nil, err
	}
	defer os.Remove(keyFile)

	args := p.buildBaseSSHArgs(keyFile, []string{
		"-L", fmt.Sprintf("%s:%s", SSHTunnelListenPort, net.JoinHostPort(p.sshArgs.RemotePostgresServer, strconv.Itoa(p.sshArgs.RemotePostgresPort))), // Local port forwarding
	})

	cmd := exec.Command("ssh", args...)
	log.Debugln("About to run ssh with args: %#v", args)

	cmd.Stdout = &b
	cmd.Stderr = &b

	if err := cmd.Start(); err != nil {
		return b.Bytes(), fmt.Errorf("starting ssh command: %w", err)
	}

	if err := p.waitForTunnel(); err != nil {
		return b.Bytes(), fmt.Errorf("waiting for tunnel: %w", err)
	}

	return b.Bytes(), nil
}

func (p *PsqlCommand) waitForTunnel() error {
	addr := net.JoinHostPort(SSHTunnelListenAdder, SSHTunnelListenPort)

	// The actual timeout set to the SSH command,
	// so I add a little more to the timeout for the port to give the SSH command time to return an error
	timeout := time.Now().Add(p.sshArgs.WaitForTunnelTimeout + 5*time.Second)
	for {
		conn, _ := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if conn != nil {
			conn.Close()
			return nil
		}

		if time.Now().After(timeout) {
			return fmt.Errorf("timeout waiting for ssh tunnel")
		}

		time.Sleep(SleepBetweenCheckingTunnel)
	}
}

func (p *PsqlCommand) getSSHExecCommand(execName string, args []string) (string, []string, error) {
	if p.sshArgs == nil || p.sshArgs.TunnelMode {
		return execName, args, nil
	}

	keyFile, err := p.writeSSHKEyFile()
	if err != nil {
		return "", nil, err
	}

	sshArgs := p.buildBaseSSHArgs(keyFile, nil)

	escapedCommand := joinCommandLines(append([]string{execName}, args...))
	sshArgs = append(sshArgs, escapedCommand)
	return "ssh", sshArgs, nil
}
