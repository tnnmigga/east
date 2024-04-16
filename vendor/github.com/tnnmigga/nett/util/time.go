package util

import "time"

func NowNs() time.Duration {
	return time.Duration(time.Now().UnixNano())
}

func NowMs() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Millisecond)
}
