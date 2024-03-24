package idgen

import (
	"east/core/conf"
	"east/core/log"
	"east/core/util"
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
	serverID := uint64(conf.ServerID())
	if serverID >= 0xFFF {
		panic("UUIDGenerater.NewID server-id must be smaller than 4096")
	}
	return idgen.timestamp<<40 | idgen.index | idgen.index<<10 | serverID
}

func newMs() uint64 {
	ms := time.Now().UnixMilli()
	return uint64(ms - 1700000000000)
}

func NewUUID() uint64 {
	return uuidgen.NewID()
}
