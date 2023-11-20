package util

import (
	"eden/core/log"
	"runtime/debug"
)

func PrintPanicStack() {
	if r := recover(); r != nil {
		log.Errorf("%v: %s", r, debug.Stack())
	}
}
