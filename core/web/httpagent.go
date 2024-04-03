package web

import (
	"context"
	"east/core/basic"
	"east/core/log"
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpAgent struct {
	*gin.Engine
	svr *http.Server
}

func NewHttpAgent() *HttpAgent {
	r := gin.New()
	agent := &HttpAgent{
		Engine: r,
		svr: &http.Server{
			Handler: r,
		},
	}
	agent.Use(func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("%v: %s", r, debug.Stack())
				ctx.String(http.StatusInternalServerError, "server panic")
			}
		}()
		ctx.Next()
	})
	return agent
}

func (agent *HttpAgent) Run(addr string) error {
	agent.svr.Addr = addr
	errChan := make(chan error, 1)
	basic.Go(func() {
		err := agent.svr.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		log.Errorf("http agent ListenAndServe return error %v", err)
		errChan <- err
	})
	time.Sleep(time.Second) // 等待1秒检测端口监听
	log.Infof("http listen and serve %s", addr)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (agent *HttpAgent) Stop(timeout ...time.Duration) error {
	waitTime := time.Minute
	if len(timeout) > 0 {
		waitTime = timeout[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	return agent.svr.Shutdown(ctx)
}
