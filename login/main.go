package main

import (
	"east/core"
	"east/core/sys"
	"east/login/modules/account"
)

func main() {
	server := core.NewServer(
		account.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
