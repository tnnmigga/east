package api

import "time"

type ITimer interface {
	Create(delay time.Duration, ctx any) uint64
	Stop(timerID uint64) bool
}
