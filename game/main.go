package main

import (
	"east/define"
	"east/game/play"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/process"
	"github.com/tnnmigga/nett/mods/mongo"
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
