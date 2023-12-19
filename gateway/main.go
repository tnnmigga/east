package main

import (
	"east/core"
	"east/core/iconf"
	"east/core/sys"
	"east/core/util"
	"east/gateway/modules/tcpagent"
)

func main() {
	iconf.LoadFromJSON(util.ReadFile("configs.jsonc"))
	server := core.NewServer(
		tcpagent.New(),
	)
	stop := server.Run()
	sys.WaitExitSignal()
	stop()
}
