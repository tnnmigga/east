package util

import (
	"io"
	"os"
)

func ReadFile(name string) []byte {
	file, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return b
}
