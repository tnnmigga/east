package main

import (
	"east/define"
	"east/game/modules/play"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/mongo"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, play.New())
	modules = append(modules, mongo.New(
		define.ModMongo,
		conf.String("mongo.uri", "mongodb://localhost:27017"),
	))
	server := nett.NewServer(modules...)
	defer server.Shutdown()
	sys.WaitExitSignal()
}
