package main

import (
	"east/core"
	"east/core/iconf"
	"east/core/sys"
	"east/core/util"
	"east/game/modules/play"
)

func main() {
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	server := core.NewServer(
		play.NewModule(),
	)
	stop := server.Run()
	sys.WaitExitSignal()
	stop()
}
