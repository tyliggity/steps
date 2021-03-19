package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
)

const outputKey = "output"

func getJsonStart(data string) (index int, jsonStr string, isArray bool) {
	for i, c := range data {
		if c == '{' {
			return i, data[i:], false
		}
		if c == '[' {
			return i, data[i:], true
		}
	}
	return -1, "", false
}

func Format(dataB []byte) (string, error) {
	ret := strings.Builder{}
	ret.WriteString(env.EndMarker())

	data := strings.TrimSpace(string(dataB))
	index, jsonStr, isArray := getJsonStart(data)
	if index < 0 {
		return "", fmt.Errorf("can't find json in given data")
	}

	_, err := gabs.ParseJSON([]byte(jsonStr))
	for err != nil {
		log.Debugln("Can't parse json in index %d, searching for another json. Err: %v", index, err)

		index = -1
		if len(jsonStr) > 2 {
			index, jsonStr, isArray = getJsonStart(jsonStr[1:])
		}
		if index < 0 {
			return "", fmt.Errorf("can't find parsable json in given data")
		}

		_, err = gabs.ParseJSON([]byte(jsonStr))
	}

	if isArray {
		// I don't use gabs to set the array in output, cause I don't want to marshal again the JSON
		// Both because performance and because I don't want it to change how the output had originally displayed
		ret.WriteString(`{"` + outputKey + `":`)
		ret.WriteString(jsonStr)
		ret.WriteRune('}')
	} else {
		ret.WriteString(jsonStr)
	}
	return ret.String(), nil
}

// Wraps the given output inside a json as value of output: {"output": "<s>"}
func FormatAsJSONOutput(s []byte) ([]byte, error) {
	result := map[string]string{outputKey: string(s)}
	return json.Marshal(result)
}
