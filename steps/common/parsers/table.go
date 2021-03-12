package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

/*
This will split a table in the following format into a array of maps

Turn this:
NAME                                                   CDS
details-v1-79c697d759-fh8qj.default                    SYNCED
istio-egressgateway-c999fffd6-xfjqn.istio-system       SYNCED

Into:
{{"NAME": "details-v1-79c697d759-fh8qj.default", "CDS": "SYNCED"}, {"NAME": "istio-egressgateway-c999fffd6-xfjqn.istio-system", "CDS": "SYNCED"}}
*/

func ParseTable(s string) ([]map[string]string, error) {
	splitRegex, err := regexp.Compile("([ ]{2,}|\\t+)")
	if err != nil {
		return nil, fmt.Errorf("failed compiling regexp: %w", err)
	}

	lines := SplitLines(s)
	if len(lines) == 0 {
		return nil, fmt.Errorf("the string doesnt seem to contain any lines")
	}

	if len(lines) == 1 {
		return []map[string]string{}, nil
	}

	headers := splitRegex.Split(lines[0], -1)
	if len(headers) == 0 {
		return nil, fmt.Errorf("no headers found")
	}

	result := make([]map[string]string, len(lines)-1)
	for i, line := range lines[1:] {
		values := splitRegex.Split(line, -1)

		if len(values) > len(headers) {
			return nil, fmt.Errorf("'%s' contains more values (%d) than headers (%d)", line, len(values), len(headers))
		}

		currentEntry := map[string]string{}

		for vi, v := range values {
			currentEntry[headers[vi]] = v
		}

		result[i] = currentEntry
	}

	return result, nil
}

func ParseTableToJSON(s string) ([]byte, error) {
	table, err := ParseTable(s)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(table)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling result to json: %w", err)
	}

	return result, nil

}
