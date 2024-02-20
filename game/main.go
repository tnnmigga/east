package main

import (
	"east/core"
	"east/core/infra/mongo"
	"east/core/sys"
	"east/game/modules/play"
)

func main() {
	server := core.NewServer(
		play.New(),
		mongo.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
