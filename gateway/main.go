package main

import (
	"east/core"
	"east/core/health"
	"east/core/iconf"
	"east/core/util"
	"east/gateway/modules/tcpagent"
)

func main() {
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	server := core.NewServer(
		tcpagent.New(),
	)
	stop := server.Run()
	health.WaitExitSignal()
	stop()
}
