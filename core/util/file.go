package util

import (
	"io"
	"log"
	"os"
)

func ReadFile(name string) []byte {
	file, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		log.Panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	return b
}
