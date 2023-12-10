package utils

import (
	"east/core/log"
	"east/core/utils/algorithms"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func RandomInterval[T constraints.Signed](low, high T) T {
	a, b := int64(low), int64(high)
	return T(rand.Int63n(b-a+1) + a)
}

func RandomIntervalN[T constraints.Signed](low, high T, num T) algorithms.Set[T] {
	maxNum := high - low + 1
	if maxNum < num {
		log.Errorf("RandomIntervalN max random num not enough")
		return nil
	}
	s := make(algorithms.Set[T], num)
	if float64(num)/float64(maxNum) > 0.75 {
		for len(s) < int(num) {
			v := RandomInterval(low, high)
			if s.Find(v) {
				continue
			}
			s.Insert(v)
		}
	} else {
		for i := low; i <= high; i++ {
			s.Insert(i)
		}
		count := high - low - num + 1
		for key := range s {
			delete(s, key)
			count--
			if count <= 0 {
				break
			}
		}
	}
	return s
}
