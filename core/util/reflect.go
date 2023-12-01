package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ReflectName(v any) string {
	return reflect.TypeOf(v).Elem().Name()
}

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
