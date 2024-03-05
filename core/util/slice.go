package util

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
