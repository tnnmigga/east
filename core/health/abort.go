package health

import (
	"east/core/log"
	"runtime"
	"sync"
)

var (
	abortChan chan error
	once      sync.Once
)

func init() {
	once.Do(func() {
		abortChan = make(chan error, 1)
	})
}

func Abort(err error) {
	select {
	case abortChan <- err:
	default:
		log.Errorf("health Abort %v", err)
	}
	runtime.Goexit()
}
