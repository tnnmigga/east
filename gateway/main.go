package main

import (
	"east/core"
	"east/core/sys"
	"east/gateway/compts/tcpagent"
)

func main() {
	server := core.NewServer(
		tcpagent.New(),
	)
	server.Init()
	server.Run()
	sys.WaitExitSignal()
	server.Stop()
	server.Close()
}
