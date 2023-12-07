package util

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
