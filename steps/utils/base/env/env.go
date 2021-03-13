package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ErrNotFound struct{ env string }

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("can't find env %s")
}

func GetEnvWithDefault(env, defaultVal string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return defaultVal
}

// Get environment variable in the format of json array or just a single element
// If the env doesn't exists, returns nil and ErrNotFound
func GetSingleOrArrayEnv(env string) ([]string, error) {
	val, ok := os.LookupEnv(env)
	if !ok {
		return nil, &ErrNotFound{env: env}
	}

	if !strings.HasPrefix(val, "[") {
		return []string{val}, nil
	}

	var output []string
	if err := json.Unmarshal([]byte(val), &output); err != nil {
		return nil, err
	}

	return output, nil
}
