package util

import "reflect"

func ReflectName(v any) string {
	return reflect.TypeOf(v).Elem().Name()
}
