package get

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stackpulse/steps-sdk-go/env"
	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	PodName       string        `env:"POD_NAME"`
	AllContainers bool          `env:"ALL_CONTAINERS" envDefault:"false"`
	ContainerName string        `env:"CONTAINER_NAME"`
	Label         string        `env:"LABEL"`
	LimitBytes    int64         `env:"LIMIT_BYTES" envDefault:"31744"` // 31KB, as environment vatiable maximum length is around 32KB
	Previous      bool          `env:"PREVIOUS" envDefault:"false"`
	Since         time.Duration `env:"SINCE" envDefault:"0s"`
	SinceTime     string        `env:"SINCE_TIME"`
	Tail          int           `env:"TAIL" envDefault:"-1"`
	Timestamps    bool          `env:"TIMESTAMPS" envDefault:"false"`
}

func (a *Args) Copy() *Args {
	return &Args{
		Args:          a.Args,
		PodName:       a.PodName,
		AllContainers: a.AllContainers,
		ContainerName: a.ContainerName,
		LimitBytes:    a.LimitBytes,
		Previous:      a.Previous,
		Since:         a.Since,
		SinceTime:     a.SinceTime,
		Tail:          a.Tail,
		Timestamps:    a.Timestamps,
	}
}

type GetLogs struct {
	Args *Args
	kctl *base2.KubectlStep
}

type LogEntry struct {
	Name string `json:"name"`
	Log  string `json:"log"`
}

func validate(args *Args) error {
	if args.SinceTime != "" {
		_, err := time.Parse(time.RFC3339, args.SinceTime)
		if err != nil {
			return fmt.Errorf("can't parse SINCE_TIME using %s layout: %w", time.RFC3339, err)
		}
	}
	return nil
}

func NewGetLogs(args *Args) (*GetLogs, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := validate(args); err != nil {
		return nil, err
	}

	return &GetLogs{
		Args: args,
		kctl: kctl,
	}, nil
}

func (l *GetLogs) getContainerNames(output string) []string {
	containersRe := regexp.MustCompile(`\[(.*?)\]`)
	containersGroups := containersRe.FindAllStringSubmatch(output, 2)
	containers := make([]string, 0)
	for _, containersGroup := range containersGroups {
		containers = append(containers, strings.Split(containersGroup[1], " ")...)
	}
	return containers
}

func (l *GetLogs) runWithArgs(args *Args) ([]byte, int, error) {
	cmdArgs := []string{"logs"}

	// Get logs by label or directly pod
	if args.Label != "" {
		cmdArgs = append(cmdArgs, "-l", args.Label)
	} else {
		cmdArgs = append(cmdArgs, args.PodName)
	}
	if args.ContainerName != "" {
		cmdArgs = append(cmdArgs, "-c", args.ContainerName)
	}
	cmdArgs = append(cmdArgs,
		fmt.Sprintf("--limit-bytes=%d", args.LimitBytes),
		fmt.Sprintf("--previous=%v", args.Previous),
		fmt.Sprintf("--since=%s", args.Since.String()),
		fmt.Sprintf("--since-time=%s", args.SinceTime),
		fmt.Sprintf("--tail=%d", args.Tail),
		fmt.Sprintf("--timestamps=%v", args.Timestamps),
	)

	return l.kctl.Execute(cmdArgs, base2.IgnoreFormat)
}

func (l *GetLogs) runForEachContainer(containers []string) ([]*LogEntry, int, error) {
	var outputs sync.Map
	var errs sync.Map
	var exitCodes sync.Map

	var wg sync.WaitGroup
	l.kctl.Debugln("Getting logs for containers: %#v", containers)
	for _, container := range containers {
		wg.Add(1)
		go func(containerName string) {
			argCopy := l.Args.Copy()
			argCopy.ContainerName = containerName
			output, exitCode, err := l.runWithArgs(argCopy)
			outputs.Store(containerName, output)
			errs.Store(containerName, err)
			exitCodes.Store(containerName, exitCode)
			wg.Done()
		}(container)
	}
	wg.Wait()

	logEntries := make([]*LogEntry, 0)
	finalExitCode := 0
	var finalErr error

	outputs.Range(func(key, value interface{}) bool {
		logEntries = append(logEntries, &LogEntry{
			Name: key.(string),
			Log:  string(value.([]byte)),
		})
		return true
	})

	exitCodes.Range(func(key, value interface{}) bool {
		if value.(int) != 0 {
			finalExitCode = value.(int)
			return false
		}
		return true
	})

	errs.Range(func(key, value interface{}) bool {
		if value != nil {
			finalErr = multierror.Append(finalErr, fmt.Errorf("can't get logs for container %s: %w", key.(string), value.(error)))
		}
		return true
	})

	return logEntries, finalExitCode, finalErr
}

func (l *GetLogs) GetRaw() (logs []*LogEntry, exitCode int, err error) {
	output, exitCode, err := l.runWithArgs(l.Args)
	if err == nil {
		return []*LogEntry{{Name: l.Args.PodName, Log: string(output)}}, exitCode, nil
	}

	if strings.Contains(string(output), "a container name must be specified for pod") {
		if l.Args.AllContainers {
			l.kctl.Debugln("Multiple containers found, running for each container")
			return l.runForEachContainer(l.getContainerNames(string(output)))
		}

		err = fmt.Errorf("must specify CONTAINER_NAME or ALL_CONTAINERS env: %w", err)
	}

	// If we are here, we have error for sure
	err = fmt.Errorf("error occured. err: %w;\nOriginal output: %s", err, string(output))
	if exitCode == 0 {
		exitCode = 1
	}
	return nil, exitCode, err
}

func (l *GetLogs) Get() (output []byte, exitCode int, err error) {
	logs, exitCode, err := l.GetRaw()
	if err != nil {
		return []byte{}, exitCode, err
	}

	outputMap := map[string][]*LogEntry{
		"logs": logs,
	}

	output, err = json.Marshal(outputMap)
	if err != nil {
		output = []byte{}
	}
	return output, exitCode, err
}

func (l *GetLogs) parseAsString(logs []*LogEntry) string {
	output := strings.Builder{}
	for _, log := range logs {
		output.WriteString(log.Name)
		output.WriteString(":\n=========================\n")
		output.WriteString(log.Log)
		output.WriteString("\n")
	}
	return output.String()
}

func (l *GetLogs) Parse(output []byte) (string, error) {
	var logsMap map[string][]*LogEntry
	if err := json.Unmarshal(output, &logsMap); err != nil {
		return "", fmt.Errorf("can't parse output as json: %w", err)
	}
	logs, ok := logsMap["logs"]
	if !ok {
		return "", fmt.Errorf("not logs key in output map")
	}

	if !env.FormatterIs(env.JsonFormat) {
		return l.parseAsString(logs), nil
	}

	var parsedBytes []byte
	var err error
	if l.Args.Pretty {
		parsedBytes, err = json.MarshalIndent(logsMap, "", "  ")
	} else {
		parsedBytes, err = json.Marshal(logsMap)
	}

	if err != nil {
		return "", fmt.Errorf("can't marshal output as json: %w", err)
	}
	return string(parsedBytes), err

}
