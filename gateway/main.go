package main

import (
	"east/core"
	"east/core/health"
	"east/core/iconf"
	"east/core/util"
	"east/game/modules/play"
)

func main() {
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	server := core.NewServer(
		play.NewModule(),
	)
	stop := server.Run()
	health.WaitExitSignal()
	stop()
}
