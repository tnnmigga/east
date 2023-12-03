package tcpagent

import (
	"east/core/message"
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
	go agent.recv()
	go agent.send()
}

func (agent *userAgent) recv() {
	sizeBuf := make([]byte, 4)
	var bufLen uint32 = 1024
	msgBuf := make([]byte, bufLen)
	_, err := io.ReadFull(agent.conn, sizeBuf)
	if err != nil {
		return
	}
	size := binary.LittleEndian.Uint32(sizeBuf)
	if size > bufLen {
		bufLen = size
		msgBuf = make([]byte, bufLen)
	}
	_, err = io.ReadFull(agent.conn, msgBuf[:size])
	if err != nil {
		return
	}
	message.Cast(agent.userInfo.serverID, &pb.C2SPackage{
		UserID: agent.userInfo.userID,
		Body:   msgBuf[:size],
	})
}

func (agent *userAgent) send() {
	for b := range agent.mq {
		agent.conn.Write(b)
	}
}

func (agent *userAgent) close() {

}