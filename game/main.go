package main

import (
	"east/core"
	"east/core/iconf"
	"east/core/log"
	"east/core/sys"
	"east/core/util"
	"east/game/modules/play"
)

func main() {
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	log.Init()
	server := core.NewServer(
		play.New(),
	)
	stop := server.Run()
	sys.WaitExitSignal()
	stop()
}
