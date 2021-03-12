package base

import (
	"fmt"
	"strconv"
	"strings"
)

// BuildPidFuncQuery transforms a pids list (e.g ["1","2","3"]") and a function (e.g. "pg_cancel_backend")
// into a SQL select query that performs the function on each pid
func BuildPidFuncQuery(pids []string, function string) (string, error) {
	var exprs []string
	for _, v := range pids {
		pid, err := strconv.Atoi(v)
		if err != nil || pid <= 0 {
			continue
		}

		exprs = append(exprs, fmt.Sprintf("%s(%d) as pid_%d", function, pid, pid))
	}

	if len(exprs) <= 0 {
		return "", fmt.Errorf("invalid pid list")
	}

	return "SELECT " + strings.Join(exprs, ","), nil
}
