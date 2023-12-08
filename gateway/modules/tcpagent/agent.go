package tcpagent

import (
	"east/core/log"
	"east/core/msgbus"
	"east/define"
	"east/pb"
	"encoding/binary"
	"io"
	"net"
)

type userInfo struct {
	userID   uint64
	serverID uint32
}

type userAgent struct {
	userInfo userInfo
	conn     net.Conn
	mq       chan []byte
}

func (agent *userAgent) run() {
	log.Infof("new agent") // 临时
	go agent.recv()
	go agent.send()
}

func (agent *userAgent) recv() {
	var sizeBuf [4]byte
	var bufLen uint32 = 1024
	msgBuf := make([]byte, bufLen)
	for {
		if _, err := io.ReadFull(agent.conn, sizeBuf[:]); err != nil {
			return
		}
		size := binary.LittleEndian.Uint32(sizeBuf[:])
		if size > bufLen {
			bufLen = size
			msgBuf = make([]byte, bufLen)
		}
		if _, err := io.ReadFull(agent.conn, msgBuf[:size]); err != nil {
			return
		}
		msgbus.Broadcast(define.ServTypGame, &pb.C2SPackage{
			UserID: agent.userInfo.userID,
			Body:   msgBuf[:size],
		})
	}
}

func (agent *userAgent) send() {
	for b := range agent.mq {
		agent.conn.Write(b)
	}
}

func (agent *userAgent) close() {

}
