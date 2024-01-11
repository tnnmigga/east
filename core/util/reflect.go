package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

func String(v any) string {
	type IString interface{ String() string }
	if v0, ok := v.(IString); ok {
		return fmt.Sprintf("[ %s: %s ]", StructName(v0), v0.String())
	}
	if b, err := json.Marshal(v); err == nil {
		return fmt.Sprintf("[ %s: %s ]", StructName(v), string(b))
	}
	return fmt.Sprintf("[ %s: %v ]", StructName(v), v)
}

func Address(fn any) uint64 {
	value := reflect.ValueOf(fn)
	ptr := unsafe.Pointer(value.Pointer())
	addr := uintptr(ptr)
	return uint64(addr)
}
