package main

import (
	"east/define"
	"east/game/modules/play"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/infra/mongo"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	server := nett.NewServer(
		play.New(),
		mongo.New(define.ModTypMongo, conf.String("mongo.uri", "mongodb://localhost:27017")),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
