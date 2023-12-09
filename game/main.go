package main

import (
	"east/core"
	"east/core/iconf"
	"east/core/sys"
	"east/core/utils"
	"east/game/modules/play"
)

func main() {
	iconf.LoadFromJSON(utils.ReadFile("configs.jsonc"))
	server := core.NewServer(
		play.New(),
	)
	stop := server.Run()
	sys.WaitExitSignal()
	stop()
}
