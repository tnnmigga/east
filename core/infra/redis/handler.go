package redis

import (
	"context"
	"east/core/sys"
	"time"
)

func (m *module) onOpGet(req *OpGet, resolve func(any), reject func(error)) {
	sys.GoWithGroup(req.Key, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		value, err := m.cli.Get(ctx,req.Key).Result()
		cancel()
		if err != nil {
			reject(err)
			return
		}
		resolve(value)
	})
}