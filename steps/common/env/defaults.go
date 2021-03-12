package env

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	ReadTimeoutEnv = "SP_READ_TIMEOUT_SECONDS"
	EndMarkerEnv   = "SP_END_MARKER"
	DebugEnv       = "SP_DEBUG"
	FormatterEnv   = "SP_FORMATTER"
)

const (
	ReadTimeoutDefault = 30 * time.Second
	EndMarkerDefault   = "<-- END -->"
	DebugDefault       = false
	FormatterDefault   = "json"
)

var (
	readTimeoutActual *time.Duration
	endMarkerActual   *string
	debugActual       *bool
	formatterActual   *string
)

func Debug() bool {
	if debugActual != nil {
		return *debugActual
	}

	debug := DebugDefault
	if val, ok := os.LookupEnv(DebugEnv); ok {
		debugParsed, err := strconv.ParseBool(val)
		if err == nil {
			debug = debugParsed
		}
	}

	debugActual = &debug
	return debug
}

func ReadTimeout() time.Duration {
	if readTimeoutActual != nil {
		return *readTimeoutActual
	}

	timeout := ReadTimeoutDefault
	if val, ok := os.LookupEnv(ReadTimeoutEnv); ok {
		seconds, err := strconv.Atoi(val)
		if err == nil {
			timeout = time.Duration(seconds) * time.Second
		}
	}
	readTimeoutActual = &timeout
	return timeout
}

func EndMarker() string {
	if endMarkerActual != nil {
		return *endMarkerActual
	}

	endMarker := GetEnvWithDefault(EndMarkerEnv, EndMarkerDefault)
	endMarkerActual = &endMarker
	return endMarker
}

func Formatter() string {
	if formatterActual != nil {
		return *formatterActual
	}

	data, err := ioutil.ReadFile(FormatOverrideFile)
	formatter := string(data)
	if err != nil {
		formatter = GetEnvWithDefault(FormatterEnv, FormatterDefault)
	}
	formatterActual = &formatter
	return formatter
}

func FormatterIs(formatter SpFormatter) bool {
	return Formatter() == string(formatter)
}
