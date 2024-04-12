package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"
)

const (
	PkgSizeByteLen   = 4
	DefaultBufferLen = 1024
)

func GetTCPBindAddress() string {
	defaultAddr := fmt.Sprintf(":%d", conf.ServerID+0x1FFE)
	return conf.String("agent.tcp.address", defaultAddr)
}

func NewTCPListener(manager *AgentManager) IListener {
	tcp := &TCPListener{
		manager: manager,
	}
	return tcp
}

type TCPListener struct {
	manager  *AgentManager
	listener net.Listener
}

func (tcp *TCPListener) Run() error {
	addr := GetTCPBindAddress()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tcp.listener = listener
	core.Go(tcp.start)
	zlog.Infof("tcp listen %s", addr)
	return nil
}

func (tcp *TCPListener) Close() {
	tcp.listener.Close()
}

func (tcp *TCPListener) start() {
	for {
		conn, err := tcp.listener.Accept()
		if err != nil {
			zlog.Warnf("tcp accept error %v", err)
			return
		}
		zlog.Debug("new conn: ", conn.RemoteAddr())
		tcpConn := NewTCPConn(conn)
		tcp.manager.OnConnect(tcpConn)
	}
}

type TCPConn struct {
	agent IAgent
	conn  net.Conn
	wBuf  []byte
}

func NewTCPConn(conn net.Conn) *TCPConn {
	return &TCPConn{
		conn: conn,
		wBuf: make([]byte, DefaultBufferLen),
	}
}

func (c *TCPConn) Write(data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	buf := c.wBuf
	pkgSize := len(data) + PkgSizeByteLen
	if pkgSize < DefaultBufferLen {
		buf = make([]byte, pkgSize)
	}
	binary.LittleEndian.PutUint32(buf, uint32(len(data)))
	copy(buf[PkgSizeByteLen:], data)
	_, err := c.conn.Write(buf[:pkgSize])
	zlog.Debugf("agent write %v", data)
	return err
}

func (c *TCPConn) ReadLoop(ctx context.Context) {
	core.Go(func() {
		sizeBuffer := make([]byte, PkgSizeByteLen)
		for {
			err := c.Read(sizeBuffer, PkgSizeByteLen)
			if util.ContextDone(ctx) {
				return
			}
			if err != nil {
				c.agent.OnReadError(err)
				continue
			}
			pkgSize := int(binary.LittleEndian.Uint32(sizeBuffer))
			if pkgSize == 0 {
				zlog.Debugf("receive ping")
				continue // 心跳包
			}
			// 每次创建一个新缓冲区
			// 防止在传递过程中可能出现的slice并发读写
			buffer := make([]byte, pkgSize)
			err = c.Read(buffer, pkgSize)
			if err != nil {
				c.agent.OnReadError(err)
				continue
			}
			c.agent.OnMessage(buffer)
		}
	})
}

func (c *TCPConn) Read(buf []byte, n int) error {
	c.conn.SetReadDeadline(time.Now().Add(MaxAliveTime))
	if _, err := io.ReadFull(c.conn, buf[:n]); err != nil {
		return err
	}
	return nil
}

func (c *TCPConn) Close() {
	c.conn.Close()
}

func (c *TCPConn) BindAgent(agent IAgent) {
	c.agent = agent
}
