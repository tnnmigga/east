package main

import (
	"east/core"
	"east/core/sys"
	"east/game/modules/play"
)

func main() {
	server := core.NewServer(
		play.New(),
	)
	server.Init()
	server.Run()
	sys.WaitExitSignal()
	server.Stop()
	server.Close()
}
