package main

import (
	"east/define"
	"east/door/modules/account"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/redis"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, account.New())
	modules = append(modules, redis.New(
		define.ModRedis,
		conf.String("redis.address", "localhost:6379"),
		conf.String("redis.username", ""),
		conf.String("redis.password", ""),
	))
	server := nett.NewServer(modules...)
	defer server.Exit()
	sys.WaitExitSignal()
}
