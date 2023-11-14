package util

import (
	"eden/common/log"
	"runtime/debug"
)

func PrintPanicStack() {
	if r := recover(); r != nil {
		log.Errorf("%v: %s", r, debug.Stack())
	}
}
