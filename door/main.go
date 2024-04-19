package main

import (
	"east/define"
	"east/door/account"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/process"
	"github.com/tnnmigga/nett/mods/mysql"
	"github.com/tnnmigga/nett/mods/redis"
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
