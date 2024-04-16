package sys

import (
	"os"
	"os/signal"
	"syscall"
)

var exitSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}

func WaitExitSignal() os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, exitSignals...)
	return <-c
}


