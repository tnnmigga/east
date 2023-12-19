package sys

import (
	"context"
	"east/core/log"
	"east/core/util"
	"sync"
	"time"
)

var (
	rootCtx, cancelGo = context.WithCancel(context.Background())
	wg                = &sync.WaitGroup{}
)

type (
	gocall interface {
		func(context.Context) | func()
	}
)

func Go[T gocall](fn T) {
	switch f := any(fn).(type) {
	case func(context.Context):
		wg.Add(1)
		go func() {
			defer util.RecoverPanic()
			defer wg.Done()
			f(rootCtx)
		}()
	case func():
		wg.Add(1)
		go func() {
			defer util.RecoverPanic()
			defer wg.Done()
			f()
		}()
	}
}

func GoWithTimeout(duration time.Duration, fn func(context.Context)) {
	wg.Add(1)
	go func() {
		defer util.RecoverPanic()
		defer wg.Done()
		ctx, _ := context.WithTimeout(rootCtx, duration)
		fn(ctx)
	}()
}

func WaitGoDone(maxWaitTime time.Duration) {
	cancelGo()
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
		log.Errorf("wait goroutine exit timeout")
	}
}
