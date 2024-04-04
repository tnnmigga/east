package main

import (
	"east/gateway/modules/agent"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	server := nett.NewServer(
		agent.New(agent.AgentTypeTCP),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
