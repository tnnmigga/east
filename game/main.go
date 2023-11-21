package main

import (
	"eden/core"
	"eden/game/modules/play"
)

func main() {
	server := core.Server{}
	server.Run(
		play.NewModule(),
	)
}
