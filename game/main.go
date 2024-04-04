package main

import (
	"east/game/modules/play"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/infra/mongo"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	server := nett.NewServer(
		play.New(),
		mongo.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
