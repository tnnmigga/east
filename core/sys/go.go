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
	wkg               = newWorkerGroup()
)

func newWorkerGroup() *workerGroup {
	return &workerGroup{
		workerPool: sync.Pool{
			New: func() any {
				return &worker{
					pending: make(chan func(), 256),
				}
			},
		},
	}
}

type gocall interface {
	func(context.Context) | func()
}

type workerGroup struct {
	group      sync.Map
	workerPool sync.Pool
	mu         sync.Mutex
}

func (wkg *workerGroup) run(key string, fn func()) {
	wkg.mu.Lock()
	var w *worker
	value, ok := wkg.group.Load(key)
	if !ok {
		w = wkg.workerPool.Get().(*worker)
		w.key = key
		wkg.group.Store(key, w)
	} else {
		w = value.(*worker)
	}
	w.count++
	pending := w.count
	wkg.mu.Unlock()
	w.pending <- fn
	if pending == 1 {
		Go(w.work)
	}
}

type worker struct {
	key     string
	pending chan func()
	count   int32
}

func (w *worker) work() {
	for {
		select {
		case fn := <-w.pending:
			util.ExecAndRecover(fn)
			w.count--
		default:
			wkg.mu.Lock()
			var empty bool
			if w.count == 0 {
				wkg.group.Delete(w.key)
				wkg.workerPool.Put(w)
				empty = true
			}
			wkg.mu.Unlock()
			if empty {
				return
			}
		}
	}
}

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

func GoWithGroup(key string, fn func()) {
	wkg.run(key, fn)
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
