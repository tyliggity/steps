package sort

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"maze.io/x/duration.v1"
	"sort"
)

func tryCompareDurations(valA, valB string) (res bool, ok bool) {
	durationA, errA := duration.ParseDuration(valA)

	durationB, errB := duration.ParseDuration(valB)

	if errA != nil && errB != nil {
		return false, false
	}

	if errA != nil {
		// If B is parsed but A is not, I want A to be "larger" then B, because I want A will be at the end of the list
		return false, true
	}

	if errB != nil {
		// If A is parsed but B is not, I want B to be "larger" then A, because I want B will be at the end of the list
		return true, true
	}

	return durationA < durationB, true
}

func tryCompareStrings(valA, valB interface{}) (res bool, ok bool) {
	valAStr, ok := valA.(string)
	if !ok {
		return false, false
	}

	valBStr, ok := valB.(string)
	if !ok {
		return false, false
	}

	// Try to parse it as duration
	if ret, ok := tryCompareDurations(valAStr, valBStr); ok {
		return ret, true
	}

	return valAStr < valBStr, true
}

func tryCompareInts(valA, valB interface{}) (res bool, ok bool) {
	valAInt, ok := valA.(int)
	if !ok {
		return false, false
	}

	valBInt, ok := valB.(int)
	if !ok {
		return false, false
	}

	return valAInt < valBInt, true
}

func tryCompareFloats(valA, valB interface{}) (res bool, ok bool) {
	valAFlt, ok := valA.(float64)
	if !ok {
		return false, false
	}

	valBFlt, ok := valB.(float64)
	if !ok {
		return false, false
	}

	return valAFlt < valBFlt, true
}

func JSON(items []*gabs.Container, jsonP string, desc bool) ([]*gabs.Container, error) {
	newItems := make([]*gabs.Container, len(items))
	for i, item := range items {
		newItem, err := gabs.ParseJSON(item.Bytes()) // The given gabs are not always parsed correctly, so I re-parse them.
		if err != nil {
			return nil, fmt.Errorf("re-parsing item %s: %w", item.String(), err)
		}
		newItems[i] = newItem
	}

	sort.Slice(newItems, func(i, j int) bool {
		if !newItems[i].ExistsP(jsonP) {
			return false
		}
		if !newItems[j].ExistsP(jsonP) {
			return true
		}

		valA := newItems[i].Path(jsonP).Data()
		valB := newItems[j].Path(jsonP).Data()

		if desc {
			valA, valB = valB, valA
		}

		if res, ok := tryCompareStrings(valA, valB); ok {
			return res
		}

		if res, ok := tryCompareInts(valA, valB); ok {
			return res
		}

		if res, ok := tryCompareFloats(valA, valB); ok {
			return res
		}

		return false
	})

	return newItems, nil
}
