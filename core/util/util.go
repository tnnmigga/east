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
