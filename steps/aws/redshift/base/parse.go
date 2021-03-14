package base

import (
	"fmt"
	"strings"
)

// ParsePidQueryResults translates the pid function query results into separated lists of pids divided by their status
func ParsePidQueryResults(results *RedshiftQueryResults) (success []string, failed []string, err error) {
	if len(results.Headers) <= 0 {
		return nil, nil, fmt.Errorf("no query results headers")
	}

	if len(results.Values) <= 0 || len(results.Values[0]) <= 0 {
		return nil, nil, fmt.Errorf("no query results")
	}

	if len(results.Values) > 1 {
		return nil, nil, fmt.Errorf("pid query result has more than one result (%d)", len(results.Values))
	}

	if len(results.Values[0]) != len(results.Headers) {
		return nil, nil, fmt.Errorf("mismatch between number of results (%d) and headers (%d)", len(results.Values[0]), len(results.Headers))
	}

	pids := make(map[string]struct{})
	success = make([]string, 0)
	failed = make([]string, 0)

	values := results.Values[0]
	for i, header := range results.Headers {
		parts := strings.Split(header, "_")
		if len(parts) < 2 {
			return nil, nil, fmt.Errorf("invalid query header '%s': %w", header, err)
		}
		pid := parts[1]

		if _, found := pids[pid]; found {
			return nil, nil, fmt.Errorf("pid '%s' appears more than once: %w", pid, err)
		} else {
			pids[pid] = struct{}{}
		}

		if values[i] == "1" {
			success = append(success, pid)
		} else {
			failed = append(failed, pid)
		}
	}

	return success, failed, nil
}
