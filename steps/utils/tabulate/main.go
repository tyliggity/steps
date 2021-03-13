package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	"sort"
	"strconv"
	"strings"
)

const HeadersByMapKeys = "by_keys"

type Args struct {
	Input           string   `env:"INPUT,required"`
	Headers         []string `env:"HEADERS" envDefault:"by_keys"`
	MapGroup        []string `env:"MAP_GROUP"`
	GroupHeaderName string   `env:"GROUP_HEADER_NAME" envDefault:"group"`
	ColumnWidth     int      `env:"COLUMN_WIDTH" envDefault:"0"`
	AutoWrap        bool     `env:"AUTO_WRAP" envDefault:"true"`
	ShowIndexes     bool     `env:"SHOW_INDEXES" envDefault:"false"`
	Markdown        bool     `env:"MARKDOWN" envDefault:"false"`
	MaxColumnLength int      `env:"MAX_COLUMN_LENGTH"  envDefault:"0"`
}

type Output struct {
	TabulateOutput string `json:"tabulate_output"`
}

type Tabulate struct {
	args         Args
	mergeByCells []int
	headers      []string
	input        [][]string
}

func (t *Tabulate) Init() error {
	var args Args
	if err := env.Parse(&args); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	t.args = args

	if unmarshalMap, ok := t.unmarshalMap(args.Input); ok {
		t.handleMapArray(args, unmarshalMap)
		return nil
	}

	if unmarshalArrays, ok := t.unmarshalArray(args.Input); ok {
		return t.handleStringArrays(args, unmarshalArrays)
	}

	if unmarshalMapArrays, ok := t.unmarshalMapArrays(args.Input); ok {
		t.handleMapOfArrays(args, unmarshalMapArrays)
		return nil
	}

	return fmt.Errorf("input should a valid JSON of map array [{...},{...}] or array of arrays [[...],[...]] or map of map arrays {\"x\":[{...},{...}], ...}")
}

func (t *Tabulate) Run() (int, []byte, error) {
	var sb strings.Builder
	table := tablewriter.NewWriter(&sb)
	table.SetRowLine(true)
	table.SetAutoWrapText(t.args.AutoWrap)

	if t.args.ColumnWidth > 0 {
		table.SetColWidth(t.args.ColumnWidth)
	}

	if t.mergeByCells != nil {
		table.SetAutoMergeCellsByColumnIndex(t.mergeByCells)
	}

	if t.args.Markdown {
		table.SetRowLine(false)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	}

	table.SetHeader(t.headers)
	table.AppendBulk(t.input)
	table.Render()

	output := &Output{TabulateOutput: sb.String()}
	res, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, res, fmt.Errorf("encode json output: %w", err)
	}

	return step.ExitCodeOK, res, nil
}

func (t *Tabulate) buildHeaders(input []map[string]string) []string {
	appendedHeaders := make(map[string]struct{})
	finalHeaders := make([]string, 0)
	for _, currentInput := range input {
		for key := range currentInput {
			if _, ok := appendedHeaders[key]; ok {
				continue
			}
			appendedHeaders[key] = struct{}{}
			finalHeaders = append(finalHeaders, key)
		}
	}
	sort.Strings(finalHeaders) // I want the user will always get the same table with the same input (and because that's map the order is not guarantee)
	return finalHeaders
}

func (t *Tabulate) buildGroupMapHeaders(input map[string][]map[string]string) []string {
	flatItems := make([]map[string]string, 0)
	for _, groupVal := range input {
		flatItems = append(flatItems, groupVal...)
	}
	return t.buildHeaders(flatItems)
}

func (t *Tabulate) formatInput(args Args, input string) string {
	input = strings.Replace(input, "\r", "", -1)
	if args.MaxColumnLength > 0 && len(input) > args.MaxColumnLength {
		input = input[:args.MaxColumnLength] + "..."
	}
	return input
}

func (t *Tabulate) inputForMapArray(args Args, group string, input []map[string]string, headers []string) [][]string {
	appendedIndex := 0
	if args.ShowIndexes {
		appendedIndex = 1
	}

	appendedGroup := 0
	if group != "" {
		appendedGroup = 1
	}

	finalInput := make([][]string, 0, len(input))
	for i, currentInput := range input {
		currentArrayInput := make([]string, 0, len(headers)+appendedIndex+appendedGroup)
		if args.ShowIndexes {
			currentArrayInput = append(currentArrayInput, strconv.Itoa(i))
		}
		if group != "" {
			currentArrayInput = append(currentArrayInput, group)
		}

		for _, header := range headers {
			currentArrayInput = append(currentArrayInput, t.formatInput(args, currentInput[header]))
		}

		finalInput = append(finalInput, currentArrayInput)
	}

	return finalInput
}

func (t *Tabulate) handleMapArray(args Args, input []map[string]string) {
	headers := args.Headers
	if len(args.Headers) == 1 && args.Headers[0] == HeadersByMapKeys {
		headers = t.buildHeaders(input)
	}

	finalInput := t.inputForMapArray(args, "", input, headers)

	if args.ShowIndexes {
		headers = append([]string{"IDX"}, headers...)
	}
	t.headers = headers
	t.input = finalInput
}

func (t *Tabulate) handleMapOfArrays(args Args, input map[string][]map[string]string) {
	mapGroups := args.MapGroup
	if len(args.MapGroup) == 0 {
		mapGroups = make([]string, 0, len(input))
		for k, _ := range input {
			mapGroups = append(mapGroups, k)
		}
		sort.Strings(mapGroups) // I want the user will always get the same table with the same input (and because that's map the order is not guarantee)
	}

	headers := args.Headers
	if len(args.Headers) == 1 && args.Headers[0] == HeadersByMapKeys {
		headers = t.buildGroupMapHeaders(input)
	}

	finalInput := make([][]string, 0)

	for _, groupKey := range mapGroups {
		finalInput = append(finalInput, t.inputForMapArray(args, groupKey, input[groupKey], headers)...)
	}

	headers = append([]string{args.GroupHeaderName}, headers...)
	t.mergeByCells = []int{0}
	if args.ShowIndexes {
		headers = append([]string{"IDX"}, headers...)
		t.mergeByCells = []int{1}
	}

	t.headers = headers
	t.input = finalInput
}

func (t *Tabulate) handleStringArrays(args Args, input [][]string) error {
	if len(args.Headers) == 1 && args.Headers[0] == HeadersByMapKeys {
		return fmt.Errorf("can't build headers when arrays list provided, please specify headers explicitly")
	}
	t.headers = args.Headers
	for i, currentInputArr := range input {
		for j, currentInput := range currentInputArr {
			input[i][j] = t.formatInput(args, currentInput)
		}
	}
	t.input = input
	return nil
}

func (t *Tabulate) unmarshalMap(input string) ([]map[string]string, bool) {
	ret := make([]map[string]string, 0)
	if err := json.Unmarshal([]byte(input), &ret); err != nil {
		log.Debugln("Unmarshal input as JSON array maps error: %v", err)
		return nil, false
	}
	return ret, true
}

func (t *Tabulate) unmarshalArray(input string) ([][]string, bool) {
	ret := make([][]string, 0)
	if err := json.Unmarshal([]byte(input), &ret); err != nil {
		log.Debugln("Unmarshal input as JSON array of arrays error: %v", err)
		return nil, false
	}
	return ret, true
}

func (t *Tabulate) unmarshalMapArrays(input string) (map[string][]map[string]string, bool) {
	ret := make(map[string][]map[string]string)
	if err := json.Unmarshal([]byte(input), &ret); err != nil {
		log.Debugln("Unmarshal input as JSON map of arrays error: %v", err)
		return nil, false
	}
	return ret, true
}

func main() {
	step.Run(&Tabulate{})
}
