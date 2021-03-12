package env

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Those parsers are from the 'env' source code
var kindParsers = map[reflect.Kind]env.ParserFunc{
	reflect.Bool: func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	},
	reflect.Int: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	},
	reflect.Int16: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	},
	reflect.Int32: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	},
	reflect.Int64: func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	},
	reflect.Int8: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	},
	reflect.Uint: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	},
	reflect.Uint16: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	},
	reflect.Uint32: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	},
	reflect.Uint64: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	},
	reflect.Uint8: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	},
	reflect.Float64: func(v string) (interface{}, error) {
		return strconv.ParseFloat(v, 64)
	},
	reflect.Float32: func(v string) (interface{}, error) {
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	},
}

// Those parsers are from the 'env' source code
var typeParsers = map[reflect.Type]env.ParserFunc{
	reflect.TypeOf(url.URL{}): func(v string) (interface{}, error) {
		u, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse URL: %v", err)
		}
		return *u, nil
	},
	reflect.TypeOf(time.Nanosecond): func(v string) (interface{}, error) {
		s, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse duration: %v", err)
		}
		return s, err
	},
}

var slicesParsers = map[reflect.Type]env.ParserFunc{
	reflect.TypeOf([]string{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]string{}))
		if err != nil {
			return nil, err
		}
		return val.([]string), nil
	},
	reflect.TypeOf([]bool{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]bool{}))
		if err != nil {
			return nil, err
		}
		return val.([]bool), nil
	},
	reflect.TypeOf([]int{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]int{}))
		if err != nil {
			return nil, err
		}
		return val.([]int), nil
	},
	reflect.TypeOf([]int8{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]int8{}))
		if err != nil {
			return nil, err
		}
		return val.([]int8), nil
	},
	reflect.TypeOf([]int16{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]int16{}))
		if err != nil {
			return nil, err
		}
		return val.([]int16), nil
	},
	reflect.TypeOf([]int32{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]int32{}))
		if err != nil {
			return nil, err
		}
		return val.([]int32), nil
	},
	reflect.TypeOf([]int64{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]int64{}))
		if err != nil {
			return nil, err
		}
		return val.([]int64), nil
	},
	reflect.TypeOf([]uint{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]uint{}))
		if err != nil {
			return nil, err
		}
		return val.([]uint), nil
	},
	reflect.TypeOf([]uint8{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]uint8{}))
		if err != nil {
			return nil, err
		}
		return val.([]uint8), nil
	},
	reflect.TypeOf([]uint16{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]uint16{}))
		if err != nil {
			return nil, err
		}
		return val.([]uint16), nil
	},
	reflect.TypeOf([]uint32{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]uint32{}))
		if err != nil {
			return nil, err
		}
		return val.([]uint32), nil
	},
	reflect.TypeOf([]uint64{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]uint64{}))
		if err != nil {
			return nil, err
		}
		return val.([]uint64), nil
	},
	reflect.TypeOf([]float32{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]float32{}))
		if err != nil {
			return nil, err
		}
		return val.([]float32), nil
	},
	reflect.TypeOf([]float64{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]float64{}))
		if err != nil {
			return nil, err
		}
		return val.([]float64), nil
	},
	reflect.TypeOf([]time.Duration{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]time.Duration{}))
		if err != nil {
			return nil, err
		}
		return val.([]time.Duration), nil
	},
	reflect.TypeOf([]url.URL{}): func(v string) (interface{}, error) {
		val, err := genericSliceParse(v, reflect.TypeOf([]url.URL{}))
		if err != nil {
			return nil, err
		}
		return val.([]url.URL), nil
	},
}

func sliceParsers(v string, typ reflect.Type) (interface{}, error) {
	if typ.Kind() == reflect.String {
		return v, nil
	}
	v = strings.TrimSpace(v)

	if parser, ok := typeParsers[typ]; ok {
		return parser(v)
	}
	if parser, ok := kindParsers[typ.Kind()]; ok {
		return parser(v)
	}

	return nil, fmt.Errorf("no parser for type %v", typ)
}

func parseJSON(jsonParsedSlice interface{}, typ reflect.Type, sliceType reflect.Type) (interface{}, error) {
	parsedSlice := reflect.ValueOf(jsonParsedSlice)
	// Because JSON parse return and array of interfaces, we should create a new array and convert the item from the json slice here after conversion
	resSlice := reflect.MakeSlice(typ, parsedSlice.Len(), parsedSlice.Cap())
	for i := 0; i < parsedSlice.Len(); i++ {
		currentVal := reflect.ValueOf(parsedSlice.Index(i).Interface())

		if parser, ok := typeParsers[sliceType]; ok {
			// If the array type is a parsable type (duration, URL) and not a built in 'Kind', parse it using given type parsers
			// for that, need to convert it to a string first
			var strValue string
			if !currentVal.Type().ConvertibleTo(reflect.TypeOf(strValue)) {
				return nil, fmt.Errorf("%q is not convertable to string and cannot be further parsed", parsedSlice.Index(i))
			}
			strValue = currentVal.Convert(reflect.TypeOf(strValue)).Interface().(string)
			// Parsing using the type parser
			val, err := parser(strValue)
			if err != nil {
				return nil, fmt.Errorf("parsing type %q of json value %q: %w", currentVal.Type(), parsedSlice.Index(i), err)
			}
			// Assigning to current val the value after parsing (with the desierd destination type)
			currentVal = reflect.ValueOf(val)
		}

		if !currentVal.Type().ConvertibleTo(sliceType) {
			return nil, fmt.Errorf("wrong type %q of json value %q", currentVal.Type(), parsedSlice.Index(i))
		}
		realVal := currentVal.Convert(sliceType)
		resSlice.Index(i).Set(realVal)
	}
	return resSlice.Interface(), nil
}

func genericSliceParse(v string, typ reflect.Type) (interface{}, error) {
	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("got non slice type")
	}
	sliceType := typ.Elem()

	jsonParsedSlice := reflect.MakeSlice(typ, 0, 0).Interface()
	// First, try unmarshal as JSON
	if err := json.Unmarshal([]byte(v), &jsonParsedSlice); err == nil && reflect.TypeOf(jsonParsedSlice).Kind() == reflect.Slice {
		return parseJSON(jsonParsedSlice, typ, sliceType)
	}

	// If that is not a JSON, try separated by ','
	vals := strings.Split(v, ",")
	resSlice := reflect.MakeSlice(typ, len(vals), len(vals))
	for i, val := range vals {
		// Parse each item depnds of it's type
		currentVal, err := sliceParsers(val, sliceType)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", val, err)
		}
		realVal := reflect.ValueOf(currentVal).Convert(sliceType)
		resSlice.Index(i).Set(realVal)
	}

	return resSlice.Interface(), nil
}

func Parse(v interface{}) error {
	return env.ParseWithFuncs(v, slicesParsers)
}
