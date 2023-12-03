package log

import (
	"fmt"

	"go.uber.org/zap"
)

func Info(a ...any) {
	zap.L().Debug("")
}

func Error(a ...any) {
	fmt.Println(a...)
	fmt.Println()
}

func Infof(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func Errorf(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func Debug(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func Debugf(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}
