package util

import (
	"reflect"
	"time"
)

func IfElse[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}

func FirstOrDefault[T any](s []T, defaultVal T) T {
	if len(s) > 0 {
		return s[0]
	}
	return defaultVal
}

func Filter[T any](items []T, fn func(T) bool) (res []T) {
	for _, item := range items {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}

func Index[T comparable](slice []T, n T) int {
	for i, v := range slice {
		if v == n {
			return i
		}
	}
	return -1
}

func Contain[T comparable](slice []T, n T) bool {
	return Index(slice, n) != -1
}

func NowNs() time.Duration {
	return time.Duration(time.Now().UnixNano())
}

func NowMs() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

// New 泛型和interface转换并去掉多层指针
func New[T any]() any {
	var v T
	typeOf := reflect.TypeOf(v)
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	return reflect.New(typeOf).Interface()
}
