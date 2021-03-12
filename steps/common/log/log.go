package log

import (
	"log"

	"github.com/stackpulse/public-steps/common/env"
)

func Log(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func Logln(msg string, args ...interface{}) {
	Log(msg+"\n", args...)
}

func Debugln(msg string, args ...interface{}) {
	if env.Debug() {
		Logln(msg, args...)
	}
}

func Debug(msg string, args ...interface{}) {
	if env.Debug() {
		Log(msg, args...)
	}
}
