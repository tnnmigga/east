package api

import "time"

type ITimer interface {
	New(delay time.Duration, ctx any) uint64
	Stop(timerID uint64) bool
}
