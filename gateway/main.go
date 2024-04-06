package main

import (
	"east/gateway/modules/agent"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/sys"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, agent.New(agent.AgentTypeTCP))
	server := nett.NewServer(modules...)
	defer server.Exit()
	sys.WaitExitSignal()
}
