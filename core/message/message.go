package message

import (
	"eden/core/configs"
	"sync"
)

var (
	once  sync.Once
	inMQ  chan<- *Package
	outMQ chan<- *Package
)

type Package struct {
	ServerID uint32
	Module   string
	Body     any
}

func Cast(serverID uint32, module string, msg any) {
	pkg := &Package{
		ServerID: serverID,
		Module:   module,
		Body:     msg,
	}
	if serverID != configs.ServerID() {
		outMQ <- &Package{}
		return
	}
	inMQ <- pkg
}

func castToInternal() {

}

// Attach 放置全局分发入口
func Attach(in, out chan<- *Package) {
	once.Do(func() {
		inMQ = in
		outMQ = out
	})
}
