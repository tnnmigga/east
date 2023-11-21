package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitStopSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
}
