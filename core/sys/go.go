package sys

import (
	"context"
	"east/core/log"
	"east/core/util"
	"sync"
	"time"
)

var (
	rootCtx, cancel = context.WithCancelCause(context.Background())
	wg              = &sync.WaitGroup{}
)

type (
	Call    func()
	WithCtx func(context.Context)
	gocall  interface {
		WithCtx | Call
	}
)

func Go[T gocall](fn T) {
	switch f := any(fn).(type) {
	case WithCtx:
		wg.Add(1)
		go func() {
			defer util.RecoverPanic()
			defer wg.Done()
			f(rootCtx)
		}()
	case Call:
		wg.Add(1)
		go func() {
			defer util.RecoverPanic()
			defer wg.Done()
			f()
		}()
	}
}

func GoWithTimeout(fn WithCtx, duration time.Duration) {
	wg.Add(1)
	go func() {
		defer util.RecoverPanic()
		defer wg.Done()
		ctx, _ := context.WithTimeout(rootCtx, duration)
		fn(ctx)
	}()
}

func WaitGoDone(maxWaitTime time.Duration) {
	c := make(chan struct{}, 1)
	timer := time.After(maxWaitTime)
	go util.ExecAndRecover(func() {
		wg.Wait()
		c <- struct{}{}
	})
	select {
	case <-c:
		return
	case <-timer:
		log.Errorf("WaitGoDone timeout")
	}
}
