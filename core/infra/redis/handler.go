package redis

import (
	"context"
	"east/core/sys"
	"time"
)

func (m *module) onOpGet(req *OpGet, resolve func(any), reject func(error)) {
	sys.GoWithGroup(req.Key, func() {
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
