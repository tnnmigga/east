package util

import (
	"io"
	"os"

	"github.com/tnnmigga/nett/zlog"
)

func ReadFile(name string) []byte {
	file, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		zlog.Panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		zlog.Panic(err)
	}
	return b
}
