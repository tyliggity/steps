package filter

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/public-steps/common/log"
	"maze.io/x/duration.v1"
	"strings"
)

type JSONFilter struct {
	Matching       map[string]string // A map between json path to the value it will be compared to. The relation between them is AND
	ComparisonFunc ComparisonFunc
}

type ComparisonFunc func(originalValue string, expectedValue string) bool

func StrEqual(orig, cmp string) bool {
	return orig == cmp
}

func StrNotEqual(orig, cmp string) bool {
	return orig != cmp
}

func StrNotContains(orig, cmp string) bool {
	return !strings.Contains(orig, cmp)
}

func DurationSmallerEqual(orig, cmp string) bool {
	originalDuration, err := duration.ParseDuration(orig)
	if err != nil {
		log.Logln("Error parsing original value '%s' as duration: %v. Excluding item from filtered.", orig, err)
		return false
	}

	sinceDuration, err := duration.ParseDuration(cmp)
	if err != nil {
		log.Logln("Error parsing expected value '%s' as duration: %v. Excluding item from filtered.", cmp, err)
		return false
	}

	return originalDuration <= sinceDuration
}

// Filter the item with AND between the filters
func isMatchFilters(item *gabs.Container, filters map[string]string, comparisonFunc ComparisonFunc) bool {
	var err error
	for filterPath, filterValue := range filters {
		item, err = gabs.ParseJSON(item.Bytes()) // For getting new pointer
		if err != nil {
			log.Logln("Reparsing item returned error: %v. Excluding item from filtered.", err)
			return false
		}
		log.Debugln("Filter for path: '%s', value: '%s'", filterPath, filterValue)
		if !item.ExistsP(filterPath) {
			log.Debugln("Path wasn't found")
			return false
		}
		val, ok := item.Path(filterPath).Data().(string)
		if !ok {
			log.Debugln("Value is not string")
			return false
		}

		if !comparisonFunc(val, filterValue) {
			log.Debugln("Comparison func returned false")
			return false
		}
	}

	return true
}

func shouldIncludeItem(item *gabs.Container, filters ...JSONFilter) bool {
	match := true
	for _, filter := range filters {
		match = match && isMatchFilters(item, filter.Matching, filter.ComparisonFunc)
	}
	return match
}

// Filtering the JSON with OR between filters, but AND between the 'Matching' field inside the filter struct
func FilterJSON(items []*gabs.Container, filters ...JSONFilter) []*gabs.Container {
	filtered := make([]*gabs.Container, 0)

	for _, item := range items {
		if shouldIncludeItem(item, filters...) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
