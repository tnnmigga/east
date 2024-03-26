package redis

import (
	"context"
	"east/core/basic"
	"time"
)

func (m *module) onOpGet(req *OpGet, resolve func(any), reject func(error)) {
	basic.GoWithGroup(req.Key, func() {
		// m.cli.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		value := m.cli.Get(ctx, req.Key).Val()
		cancel()
		// if err != nil {
		// 	reject(err)
		// 	return
		// }
		resolve(value)
	})
}
