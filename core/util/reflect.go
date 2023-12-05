package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ReflectName(v any) string {
	mType := reflect.TypeOf(v)
	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	return mType.Name()
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
