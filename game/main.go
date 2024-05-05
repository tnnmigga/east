package main

import (
	"east/define"
	"east/game/play"

	"github.com/tnnmigga/core"
	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/infra/process"
	"github.com/tnnmigga/core/mods/mongo"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, play.New())
	modules = append(modules, mongo.New(
		define.ModMongo,
		conf.String("mongo.uri", "mongodb://localhost:27017"),
		conf.String("mongo.db", "game"),
	))
	server := nett.NewServer(modules...)
	defer server.Shutdown()
	process.WaitExitSignal()
}
