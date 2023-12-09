package main

import (
	"east/core"
	"east/core/iconf"
	"east/core/sys"
	"east/core/utils"
	"east/gateway/modules/tcpagent"
)

func main() {
	iconf.LoadFromJSON(utils.ReadFile("configs.jsonc"))
	server := core.NewServer(
		tcpagent.New(),
	)
	stop := server.Run()
	sys.WaitExitSignal()
	stop()
}
