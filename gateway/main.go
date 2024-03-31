package main

import (
	"east/core"
	"east/core/sys"
	"east/gateway/modules/agent"
)

func main() {
	server := core.NewServer(
		agent.New(agent.AgentTypeTCP),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
