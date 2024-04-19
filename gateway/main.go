package main

import (
	"east/gateway/agent"

	"github.com/tnnmigga/nett"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/process"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, agent.New(conf.String("agent.type", agent.AgentTypeTCP)))
	server := nett.NewServer(modules...)
	defer server.Shutdown()
	process.WaitExitSignal()
}
