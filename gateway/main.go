package main

import (
	"east/gateway/agent"

	"github.com/tnnmigga/core"
	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/infra/process"
)

func main() {
	var modules []idef.IModule
	modules = append(modules, agent.New(conf.String("agent.type", agent.AgentTypeTCP)))
	server := nett.NewServer(modules...)
	defer server.Shutdown()
	process.WaitExitSignal()
}
