package main

import (
	"east/define"
	"east/door/account"

	"github.com/tnnmigga/core"
	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/infra/process"
	"github.com/tnnmigga/core/mods/mysql"
	"github.com/tnnmigga/core/mods/redis"
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
	modules = append(modules, mysql.New(
		define.ModMysql,
		conf.String("mysql.dsn", ""),
	))
	server := nett.NewServer(modules...)
	defer server.Shutdown()
	process.WaitExitSignal()
}
