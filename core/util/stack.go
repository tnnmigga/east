package util

import (
	"eden/core/log"
	"runtime"
	"runtime/debug"
)

func RecoverPanic() {
	if r := recover(); r != nil {
		log.Errorf("%v: %s", r, debug.Stack())
	}
}

func ExecAndRecover(fn func()) {
	defer RecoverPanic()
	fn()
}

func Caller(skip ...int) string {
	n := 1
	if len(skip) > 0 {
		n = skip[0]
	}
	pc, _, _, ok := runtime.Caller(n)
	if !ok {
		return "runtime.Caller() failed"
	}
	return runtime.FuncForPC(pc).Name()
}
