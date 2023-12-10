package util

import (
	"east/core/log"
	"east/core/util/algorithm"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func RandomInterval[T constraints.Signed](low, high T) T {
	a, b := int64(low), int64(high)
	return T(rand.Int63n(b-a+1) + a)
}

func RandomIntervalN[T constraints.Signed](low, high T, num T) algorithm.Set[T] {
	maxNum := high - low + 1
	if maxNum < num {
		log.Errorf("RandomIntervalN max random num not enough")
		return nil
	}
	var set algorithm.Set[T]
	if float64(num)/float64(maxNum) < 0.75 {
		set = make(algorithm.Set[T], num)
		for len(set) < int(num) {
			v := RandomInterval(low, high)
			if set.Find(v) {
				continue
			}
			set.Insert(v)
		}
		return set
	}
	set = make(algorithm.Set[T], maxNum)
	for i := low; i <= high; i++ {
		set.Insert(i)
	}
	count := high - low - num + 1
	for key := range set {
		delete(set, key)
		count--
		if count <= 0 {
			break
		}
	}
	return set
}
