package base

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/filter"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/sort"
	"maze.io/x/duration.v1"
)

const ItemsArrayName = "items"

type JsonParseConfig struct {
	ParseFunc JsonParseFunc
	Args      []string
}

type JsonParseFunc func(args []string, outputItem *gabs.Container) (interface{}, error)

// Returning string value of the given json path
// Args:
// 0 - The path to the string inside the JSON
// 1 (optional) - Default value in case the path wasn't found or it's not a string
// If the path wasn't found (or it's not string) will return default value if given or empty string
func JsonPathStringParser(args []string, item *gabs.Container) (interface{}, error) {
	path := args[0]
	defaultValue := ""
	if len(args) > 1 {
		defaultValue = args[2]
	}

	if !item.ExistsP(path) {
		return defaultValue, nil
	}

	val, ok := item.Path(args[0]).Data().(string)
	if !ok {
		return defaultValue, nil
	}
	return val, nil
}

// Returning JsonArray as []string value of the given json path object keys
// Args:
// 0 - The path to the object inside the JSON
// Will return empty array if the value not found or it's not an object
// Example:
// The JSON: {"items" {"foo":"bar", "lorem":[1,2,3], "ipsum": {}}
// With args: []string{"items"}
// Will return: []string{"foo", "lorem", "ipsum"}
func JsonPathObjectKeys(args []string, item *gabs.Container) (interface{}, error) {
	children := item.Path(args[0]).ChildrenMap()
	keys := make([]string, len(children))

	i := 0
	for key := range children {
		keys[i] = key
		i++
	}
	return keys, nil
}

// Returning JsonArray as []interface{} value of all the values of given key inside given array path
// Args:
// 0 - The path to the array of objects inside the JSON
// 1 - The key to extract the value from
// Will return empty array if the value not found or it's not an array
// Example:
// The JSON: {"items": [ {"foo":"bar", "foo":[1,2,3], "foo": {}, "bar": "excluded" ] }
// With args: []string{"items", "foo"}
// Will return: []string{"bar", "[1,2,3]", "{}"}
func JsonPathObjectArrayKeyValue(args []string, item *gabs.Container) (interface{}, error) {
	arrayPath := args[0]
	key := args[1]

	children := item.Path(arrayPath).Children()
	values := make([]interface{}, 0, len(children))

	for _, child := range children {
		if child.Exists(key) {
			values = append(values, child.S(key))
		}
	}
	return values, nil
}

// Returning length of JSON array
// Args:
// 0 - The path to the array inside the JSON
func JsonPathArrayLength(args []string, item *gabs.Container) (interface{}, error) {
	return len(item.Path(args[0]).Children()), nil
}

// Return Duration representation of time.Now() subtracting given json path timestamp
// Args:
// 0 - The path to the timestamp inside the JSON
// 1 - The layout of the timestamp for parsing
// Returns 1w2d3h4m representation format
func JsonPathDurationFromDate(args []string, item *gabs.Container) (interface{}, error) {
	timestampPath := args[0]
	layout := args[1]

	val, ok := item.Path(timestampPath).Data().(string)
	if !ok {
		return time.Duration(0).String(), nil
	}

	parsedTime, err := time.Parse(layout, val)
	if err != nil {
		return nil, err
	}
	// Golang default duration return just hour reorientation, this duration represents also weeks and days
	return duration.Duration(time.Now().Sub(parsedTime)).String(), nil
}

// Searching for a particular key with particular value in object array and get a value of the object if matched
// Args:
// 0 - Path to the array inside the JSON
// 1 - The key of each json object inside the array to search in
// 2 - The value we should search for
// 3 - The key of the json object inside the array to get the return value from in case the search criteria matched
// 4 - Default value in case we couldn't find the key with the given search value or the value key was empty
// 5 (optional) - If "true" returns error instead of the default value in case we couldn't find the key with the given search value or the value key was empty
// Example:
// The JSON: {"items": [ {"type":"foo", "data":"bar"}, {"type":"foo2", "data":"bar2"} ] }
// With the following args: []string{"items", "type", "foo2", "data", "DefaultVal"}
// Will return: "bar2"
// If "foo2" will not be found (or the "data" key will be empty) the return value will be "DefaultVal"
// You can add a sixth argument with "true" as value if you want the function to return error instead of returning default value
func JsonPathSearchInObjectArray(args []string, item *gabs.Container) (interface{}, error) {
	arrayPath := args[0]
	searchKey := args[1]
	searchVal := args[2]
	valueKey := args[3]
	defaultValue := args[4]
	retErrNotFound := false
	if len(args) > 5 {
		retErrNotFound = args[5] == "true"
	}

	for _, val := range item.Path(arrayPath).Children() {
		if v, _ := val.Path(searchKey).Data().(string); v == searchVal {
			foundValue, _ := val.Path(valueKey).Data().(string)
			if foundValue != "" {
				return foundValue, nil
			}
		}
	}
	if retErrNotFound {
		return nil, fmt.Errorf("key %s bot found or it's value is empty in path %s", searchKey, arrayPath)
	}
	return defaultValue, nil
}

func (k *KubectlStep) ParseOutput(output []byte, parsingConfiguration map[string]*JsonParseConfig, extraFilters ...filter.JSONFilter) (string, error) {
	parsed, err := k.parseOutput(output, parsingConfiguration, extraFilters...)
	if err != nil {
		env.SetFormatter(env.PrintFormat, true)
		return string(output), err
	}
	return parsed, nil
}

func (k *KubectlStep) parseSingleItem(item *gabs.Container, parsingConfiguration map[string]*JsonParseConfig) (*gabs.Container, error) {
	if parsingConfiguration == nil {
		return item, nil
	}

	retItem := gabs.New()
	for fieldName, parseConfig := range parsingConfiguration {
		val, err := parseConfig.ParseFunc(parseConfig.Args, item)
		if err != nil {
			return nil, err
		}
		retItem.Set(val, fieldName)
	}
	return retItem, nil
}

func (k *KubectlStep) parseItems(items []*gabs.Container, parsingConfiguration map[string]*JsonParseConfig) (*gabs.Container, error) {
	retGc := gabs.New()
	retGc.Array(ItemsArrayName)
	for _, item := range items {
		currentRetItem, err := k.parseSingleItem(item, parsingConfiguration)
		if err != nil {
			return nil, err
		}
		retGc.ArrayAppend(currentRetItem, ItemsArrayName)
	}
	return retGc, nil
}

func (k *KubectlStep) FilterResult(gc *gabs.Container, extraFilters ...filter.JSONFilter) []*gabs.Container {
	if len(extraFilters) == 0 &&
		len(k.StepArgs.FilterEqualsParsed) == 0 &&
		len(k.StepArgs.FilterContainsParsed) == 0 &&
		len(k.StepArgs.FilterNotEqualsParsed) == 0 &&
		len(k.StepArgs.FilterNotContainsParsed) == 0 {
		return nil
	}

	k.Debugln("Filtering results")
	filters := []filter.JSONFilter{
		{Matching: k.StepArgs.FilterEqualsParsed, ComparisonFunc: filter.StrEqual},
		{Matching: k.StepArgs.FilterContainsParsed, ComparisonFunc: strings.Contains},
		{Matching: k.StepArgs.FilterNotEqualsParsed, ComparisonFunc: filter.StrNotEqual},
		{Matching: k.StepArgs.FilterNotContainsParsed, ComparisonFunc: filter.StrNotContains},
	}
	filters = append(filters, extraFilters...)

	return filter.FilterJSON(gc.Path(ItemsArrayName).Children(), filters...)
}

func (k *KubectlStep) orderItems(gc *gabs.Container) *gabs.Container {
	items := gc.S(ItemsArrayName).Children()
	sorted, err := sort.JSON(items, k.StepArgs.OrderBy, k.StepArgs.OrderByDescending)
	if err != nil {
		log.Logln("Failed sorting items: %v", err.Error())
		return gc
	}

	retGc := gabs.New()
	retGc.Array(ItemsArrayName)
	for _, item := range sorted {
		retGc.ArrayAppend(item, ItemsArrayName)
	}

	return retGc
}

func (k *KubectlStep) parseOutput(output []byte, parsingConfiguration map[string]*JsonParseConfig, extraFilters ...filter.JSONFilter) (string, error) {
	if k.StepArgs.Format != "json" {
		return string(output), nil
	}

	gc, err := gabs.ParseJSON(output)
	if err != nil {
		return "", fmt.Errorf("failed parse output as json: %w", err)
	}

	var retGc *gabs.Container

	switch kind, _ := gc.Path("kind").Data().(string); kind {
	case "List":
		retGc, err = k.parseItems(gc.S("items").Children(), parsingConfiguration)
	default:
		retGc, err = k.parseItems([]*gabs.Container{gc}, parsingConfiguration)
	}

	if err != nil {
		return "", err
	}

	if k.StepArgs.OrderBy != "" {
		k.Debugln("Ordering items by %s", k.StepArgs.OrderBy)
		retGc = k.orderItems(retGc)
	}

	filtered := k.FilterResult(retGc, extraFilters...)
	if filtered != nil {
		retGc.Set(filtered, "filtered")
	}

	if k.StepArgs.Pretty {
		return retGc.StringIndent("", "  "), nil
	}
	return retGc.String(), nil
}
