package idgen

import (
	"eden/core/conf"
	"eden/core/log"
	"eden/core/util"
	"sync"
	"time"
)

var uuidgen UUIDGenerater

type UUIDGenerater struct {
	sync.Mutex
	timestamp uint64
	index     uint64
}

func (idgen *UUIDGenerater) NewID() uint64 {
	idgen.Lock()
	defer idgen.Unlock()
	ms := newMs()
	if ms != idgen.timestamp {
		idgen.timestamp = ms
		idgen.index = 0
	}
	idgen.index++
	if idgen.index > 0x3FF {
		log.Errorf("idgen uuid index over limit, caller %v", util.Caller())
		return 0
	}
	return idgen.timestamp<<40 | idgen.index | idgen.index<<10 | conf.ServerID()
}

func newMs() uint64 {
	ms := time.Now().UnixMilli()
	return uint64(ms - 1700000000000)
}

func NewUUID() uint64 {
	return uuidgen.index
}
