package main

import (
	"east/core"
	"east/core/infra/redis"
	"east/core/sys"
	"east/door/modules/account"
)

func main() {
	server := core.NewServer(
		account.New(),
		redis.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
