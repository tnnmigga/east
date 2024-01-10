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
	defer server.Exit()
	sys.WaitExitSignal()
}
