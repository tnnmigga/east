package log

import "fmt"

func Info(a ...any) {
	fmt.Println(a...)
}

func Error(a ...any) {
	fmt.Println(a...)
}

func Infof(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Errorf(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Debug(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Debugf(format string, a ...any) {
	fmt.Printf(format, a...)
}
