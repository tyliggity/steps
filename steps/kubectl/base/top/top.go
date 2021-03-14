package top

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/public-steps/common/env"
	base2 "github.com/stackpulse/public-steps/kubectl/base"
	"strings"
)

type Args struct {
	base2.Args
	ResourceType string `env:"RESOURCE_TYPE,required"`
	SortBy       string `env:"SORT_BY"`
	ResourceName string `env:"RESOURCE_NAME"`
}

var ValidResourceType = []string{"node", "nodes", "pod", "pods"}
var ValidSortBy = []string{"cpu", "memory"}

func stringInSlice(val string, items []string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}
	return false
}
func (a *Args) validate() error {
	if !stringInSlice(a.ResourceType, ValidResourceType) {
		return fmt.Errorf("RESOURCE_TYPE must be one of: %+v", ValidResourceType)
	}

	if a.SortBy != "" {
		if !stringInSlice(a.SortBy, ValidSortBy) {
			return fmt.Errorf("SORT_BY must be one of: %+v", ValidSortBy)
		}
	}
	return nil
}

type TopEntry map[string]string

type Top struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewTop(args *Args) (*Top, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := args.validate(); err != nil {
		return nil, err
	}

	return &Top{
		Args: args,
		kctl: kctl,
	}, nil
}

func (t *Top) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"top", t.Args.ResourceType, fmt.Sprintf("--sort-by=%s", t.Args.SortBy)}
	if t.Args.ResourceName != "" {
		cmdArgs = append(cmdArgs, t.Args.ResourceName)
	}

	return t.kctl.Execute(cmdArgs, base2.IgnoreFormat, base2.IgnoreFieldSelector)
}

func getHeaders(line string) []string {
	rawHeaders := strings.Fields(line)
	for i, header := range rawHeaders {
		header = strings.ToLower(header)
		if idx := strings.Index(header, "("); idx > 0 {
			description := ""
			if idxClose := strings.Index(header[idx:], ")"); idxClose > 0 {
				// Getting string inside brackets with first letter upper cased
				description = strings.Title(header[idx+1 : idx+idxClose])
			}
			header = header[:idx] + description
		}
		if idx := strings.Index(header, "%"); idx > 0 {
			header = header[:idx] + "Usage"
		}
		rawHeaders[i] = header
	}

	return rawHeaders
}

func alignFieldsWithHeaders(headers, fields []string) ([]string, bool) {
	if len(fields) == len(headers) {
		return fields, true
	}

	if len(fields) > len(headers) {
		fields = fields[:len(headers)]
		return fields, false
	}

	// Appending empty strings to align the fields with the headers
	for i := len(fields); i < len(headers); i++ {
		fields = append(fields, "")
	}
	return fields, false
}

func (t *Top) Parse(output []byte) (string, error) {
	if !env.FormatterIs(env.JsonFormat) {
		return string(output), nil
	}

	data := strings.TrimSpace(string(output))
	lines := strings.Split(data, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("no lines in data")
	}
	// First, getting the headers
	headers := getHeaders(lines[0])

	lines = lines[1:]
	items := make([]TopEntry, 0, len(lines))
	for i, line := range lines {
		currentItem := make(TopEntry)
		fields, ok := alignFieldsWithHeaders(headers, strings.Fields(line))
		if !ok {
			t.kctl.Debugln("Fields are not aligned with headers in line %d", i+1)
			t.kctl.Debugln("Line: %s; Headers: %+v", line, headers)
		}
		for j, header := range headers {
			currentItem[header] = fields[j]
		}
		items = append(items, currentItem)
	}

	gc := gabs.New()
	gc.Set(items, "items")

	if t.Args.Pretty {
		return gc.StringIndent("", "  "), nil
	}
	return gc.String(), nil
}
