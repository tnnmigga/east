package test

import (
	"sync"
	"testing"
)

const maxN = 25

var m map[int]int
var l []int
var once sync.Once

func Init() {
	once.Do(func() {
		m = map[int]int{}
		for i := 0; i < maxN; i++ {
			m[i] = i
			l = append(l, i)
		}
	})

}
func BenchmarkMap(t *testing.B) {
	Init()
	for i := 0; i < maxN; i++ {
		m[i]++
	}
}

func BenchmarkFind(t *testing.B) {
	Init()
	for i := 0; i < maxN; i++ {
		for j := 0; j < maxN; j++ {
			if l[j] == 25 {
				l[j] = 25
				break
			}
		}
	}
}
