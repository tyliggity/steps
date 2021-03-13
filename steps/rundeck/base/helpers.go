package base

import (
	"encoding/json"
	"regexp"
	"strings"
)

func jsonKeysCaseSerializer(apiOutput []byte) ([]byte, error) {
	var untypedJson interface{}
	err := json.Unmarshal(apiOutput, &untypedJson)
	if err != nil {
		return apiOutput, err
	}

	switch untypedJson.(type) {
	case map[string]interface{}:
		untypedJson = parseMap(untypedJson.(map[string]interface{}))
	case []interface{}:
		untypedJson = parseArray(untypedJson.([]interface{}))
	default:
	}

	return json.Marshal(untypedJson)
}

func parseMap(aMap map[string]interface{}) interface{} {
	for key, val := range aMap {
		if key == "options" { // options of job execution. will contain job specific options.
			continue
		}

		aMap, serializedKey := replaceKey(aMap, key)

		switch val.(type) {
		case map[string]interface{}:
			aMap[serializedKey] = parseMap(val.(map[string]interface{}))
		case []interface{}:
			aMap[serializedKey] = parseArray(val.([]interface{}))
		default:
			continue
		}
	}
	return aMap
}

func parseArray(anArray []interface{}) interface{} {
	for i, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			anArray[i] = parseMap(val.(map[string]interface{}))
		case []interface{}:
			anArray[i] = parseArray(val.([]interface{}))
		default:
			continue
		}
	}

	return anArray
}

func replaceKey(jsonMap map[string]interface{}, key string) (map[string]interface{}, string) {
	serializedKey := key
	serializedKey = replaceChars(serializedKey)

	if serializedKey != key {
		jsonMap[serializedKey] = jsonMap[key]
		delete(jsonMap, key)
	}

	return jsonMap, serializedKey
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchScore    = regexp.MustCompile("([a-z0-9])-([a-z])")
)

func replaceChars(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchScore.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
