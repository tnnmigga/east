package main

import (
	"east/core"
	"east/core/sys"
	"east/gateway/modules/tcpagent"
)

func main() {
	server := core.NewServer(
		tcpagent.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
