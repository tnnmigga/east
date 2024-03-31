package main

import (
	"east/core"
	"east/core/sys"
	"east/door/modules/account"
)

func main() {
	server := core.NewServer(
		account.New(),
	)
	defer server.Exit()
	sys.WaitExitSignal()
}
