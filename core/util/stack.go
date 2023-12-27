package util

import (
	"east/core/log"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

func Panic(err error) {
	panic(err)
}

func RecoverPanic() {
	if r := recover(); r != nil {
		log.Errorf("%v: %s", r, debug.Stack())
	}
}

func ExecAndRecover(fn func()) {
	defer RecoverPanic()
	fn()
}

// Caller 获取调用者
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

// Marker 调用位置唯一标记
// func Marker(skip ...int) string {
// 	n := 1
// 	if len(skip) > 0 {
// 		n = skip[0]
// 	}
// 	pc, codePath, codeLine, ok := runtime.Caller(n)
// 	if !ok {
// 		return ""
// 	}
// 	name := runtime.FuncForPC(pc).Name()
// 	return fmt.Sprintf("%s:%d %s", codePath, codeLine, name)
// }

// 获取包名
func PkgName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	return strings.Split(name, ".")[0]
}

// StructName 获取结构体名称
func StructName(v any) string {
	mType := reflect.TypeOf(v)
	for mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	return mType.Name()
}

// FuncName 获取函数名称
func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
