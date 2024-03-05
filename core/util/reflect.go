package util

import (
	"encoding/json"
	"fmt"
	"reflect"
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

// 泛型和interface转换并去掉多层指针
func New[T any]() any {
	var v T
	typeOf := reflect.TypeOf(v)
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	return reflect.New(typeOf).Interface()
}
