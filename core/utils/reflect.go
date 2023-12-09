package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

func String(v any) string {
	type IString interface{ String() string }
	if v0, ok := v.(IString); ok {
		return v0.String()
	}
	if b, err := json.Marshal(v); err == nil {
		return string(b)
	}
	return fmt.Sprint(v)
}

func Address(fn any) uint64 {
	value := reflect.ValueOf(fn)
	ptr := unsafe.Pointer(value.Pointer())
	addr := uintptr(ptr)
	return uint64(addr)
}
