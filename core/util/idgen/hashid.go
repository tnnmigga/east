package idgen

import (
	"crypto/md5"
	"encoding/binary"
)

func Hashuint64(s string) uint64 {
	sum := md5.Sum([]byte(s))
	return binary.LittleEndian.Uint64(sum[:])
}

func Hashuint32(s string) uint32 {
	sum := md5.Sum([]byte(s))
	return binary.LittleEndian.Uint32(sum[:])
}

func Hashint64(s string) int64 {
	sum := md5.Sum([]byte(s))
	sum[7] = 0
	return int64(binary.LittleEndian.Uint64(sum[:]))
}

func Hashint32(s string) int32 {
	sum := md5.Sum([]byte(s))
	sum[3] = 0
	return int32(binary.LittleEndian.Uint32(sum[:]))
}
