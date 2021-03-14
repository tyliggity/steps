package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type SpFormatter string

const (
	JsonFormat  SpFormatter = "json"
	RawFormat   SpFormatter = "raw"
	PrintFormat SpFormatter = "print"
)

// This file used to override SP_FORMATTER env (in case the step want to override it, because it can't modify the actual env)
const FormatOverrideFile = "/tmp/format.txt"

func GetEnvWithDefault(env, defaultVal string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return defaultVal
}

func SetFormatter(format SpFormatter, force bool) {
	_, ok := os.LookupEnv(FormatterEnv)

	formatStr := string(format)
	if force || !ok {
		os.Setenv(FormatterEnv, formatStr)
		ioutil.WriteFile(FormatOverrideFile, []byte(formatStr), os.FileMode(0644))
		formatterActual = &formatStr
	}
}

func ToTypeAwareKeyValue(v map[string]string) map[string]interface{} {
	keyVal := make(map[string]interface{})

	for key, val := range v {
		if val == "true" {
			keyVal[key] = true
			continue
		}

		if val == "false" {
			keyVal[key] = false
			continue
		}

		i, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			keyVal[key] = i
			continue
		}

		keyVal[key] = val
	}

	return keyVal
}

func ParseKeyValueEnv(v string) (interface{}, error) {
	details := map[string]string{}
	detailsPairs := strings.Split(v, ",")

	for _, pair := range detailsPairs {
		if pair == "" {
			continue
		}

		splitPair := strings.Split(pair, "=")
		if len(splitPair) != 2 {
			return nil, fmt.Errorf("invalid details pair: '%s' seperator is '=', example: 'key=value'", pair)
		}

		details[splitPair[0]] = splitPair[1]
	}

	return details, nil
}

// This type is useful in case of json input of items array
// You can find usage example in ssh-parallel command
type JSONItemsArray struct {
	Items []string
}

func ParseJSONItemsArray(v string) (interface{}, error) {
	var jsonItemsArray JSONItemsArray
	err := json.Unmarshal([]byte(v), &jsonItemsArray.Items)
	if err != nil {
		return nil, err
	}

	return jsonItemsArray, nil
}
