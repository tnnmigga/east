package redis

import (
	"context"
	"east/core/core"
	"east/core/msgbus"
	"east/core/util"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrInvalidCmd = errors.New("invalid command")

func (m *module) initHandler() {
	msgbus.RegisterRPC(m, m.onExec)
	msgbus.RegisterRPC(m, m.onExecMulti)
}

func (m *module) onExec(req *Exec, resolve func(any), reject func(error)) {
	if len(req.Cmd) < 2 {
		reject(ErrInvalidCmd)
		return
	}
	key := req.Key
	if key == "" {
		key = req.Cmd[1].(string)
	}
	core.GoWithGroup(key, func() {
		ctx, cancel := context.WithTimeout(context.Background(), util.IfElse(req.Timeout > 0, req.Timeout, 3*time.Second))
		defer cancel()
		cmd := m.cli.Do(ctx, req.Cmd...)
		result, err := cmd.Result()
		if err != nil {
			reject(err)
			return
		}
		resolve(result)
	})
}

func (m *module) onExecMulti(req *ExecMulti, resolve func(any), reject func(error)) {
	if len(req.Cmds) == 0 {
		reject(ErrInvalidCmd)
		return
	}
	core.GoWithGroup(req.Key, func() {
		ctx, cancel := context.WithTimeout(context.Background(), util.IfElse(req.Timeout > 0, req.Timeout, 3*time.Second))
		defer cancel()
		pipe := m.cli.TxPipeline()
		cmders := make([]*redis.Cmd, 0, len(req.Cmds))
		for _, cmd := range req.Cmds {
			cmders = append(cmders, pipe.Do(ctx, cmd...))
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			reject(err)
			return
		}
		results := make([]interface{}, 0, len(cmders))
		for _, cmder := range cmders {
			result, _ := cmder.Result()
			results = append(results, result)
		}
		resolve(results)
	})
}
