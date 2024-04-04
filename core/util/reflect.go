package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

// 泛型和interface转换并去掉多层指针
func New[T any]() any {
	var value T
	typeOf := reflect.TypeOf(value)
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	return reflect.New(typeOf).Interface()
}

func String(v any) string {
	type IString interface{ String() string }
	if v0, ok := v.(IString); ok {
		return fmt.Sprintf("[ %s: %s ]", TypeName(v0), v0.String())
	}
	if b, err := json.Marshal(v); err == nil {
		return fmt.Sprintf("[ %s: %s ]", TypeName(v), string(b))
	}
	return fmt.Sprintf("[ %s: %v ]", TypeName(v), v)
}

func Integer[T constraints.Integer](value any) T {
	switch v := value.(type) {
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return T(n)
	case int:
		return T(v)
	case int16:
		return T(v)
	case int32:
		return T(v)
	case int64:
		return T(v)
	case uint:
		return T(v)
	case uint16:
		return T(v)
	case uint32:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	case uintptr:
		return T(v)
	case byte:
		return T(v)
	default:
		panic(fmt.Errorf("Integer type error %v", reflect.TypeOf(value)))
	}
}
