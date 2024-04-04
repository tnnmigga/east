package main

import (
	"east/door/modules/account"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/infra/redis"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	server := nett.NewServer(
		account.New(),
		redis.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
