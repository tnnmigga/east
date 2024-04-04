package agent

import (
	"context"
	"east/core/core"
	"east/core/util"
	"east/core/zlog"
	"encoding/binary"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const (
	MaxTcpWaitTime = 5 * time.Second
	MaxTcpPkgSize  = 1024
	PkgSizeByteLen = 4
)

func NewTCPListener(manager *AgentManager, addr string) IListener {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zlog.Fatalf("tcpagent listen error %v", err)
	}
	tcp := &TCPListener{
		manager:  manager,
		listener: listener,
	}
	zlog.Infof("tcp listen %s", addr)
	return tcp
}

type TCPListener struct {
	manager  *AgentManager
	listener net.Listener
}

func (tcp *TCPListener) Run() {
	core.Go(tcp.accept)
	core.Go(func(ctx context.Context) {
		timer := time.NewTicker(MaxTcpWaitTime)
		for {
			select {
			case <-timer.C:
				tcp.killDeadAgent()
			case <-ctx.Done():
				return
			}
		}
	})
}

func (tcp *TCPListener) Close() {
	tcp.listener.Close()
	for _, agent := range tcp.manager.agents {
		agent.conn.Close()
	}
}

func (tcp *TCPListener) killDeadAgent() {
	tcp.manager.rw.RLock()
	nowNs := util.NowNs()
	for uid, agent := range tcp.manager.agents {
		state := atomic.LoadInt32(&agent.state)
		if state == AgentStateDead {
			delete(tcp.manager.agents, uid)
			continue
		}
		if state != AgentStateWait {
			continue
		}
		if agent.waitNs+MaxTcpWaitTime > nowNs {
			continue
		}
		if !atomic.CompareAndSwapInt32(&agent.state, AgentStateWait, AgentStateDead) {
			continue
		}
		delete(tcp.manager.agents, uid)
	}
	tcp.manager.rw.RUnlock()
}

func (tcp *TCPListener) accept() {
	defer util.RecoverPanic()
	for {
		conn, err := tcp.listener.Accept()
		if err != nil {
			zlog.Warnf("tcp accept error %v", err)
			return
		}
		zlog.Debug("new conn: ", conn.RemoteAddr())
		tcpConn := &TCPConn{
			conn: conn,
		}
		tcp.manager.OnConnect(tcpConn)
	}
}

type TCPConn struct {
	agent IAgent
	conn  net.Conn
}

func (c *TCPConn) Write(data []byte) error {
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err := c.conn.Write(data)
	return err
}

func (c *TCPConn) Run(ctx context.Context) {
	core.Go(c.readLoop)
}

func (c *TCPConn) readLoop() {
	buffer := make([]byte, MaxTcpPkgSize)
	for {
		err := c.Read(buffer, PkgSizeByteLen)
		if err != nil {
			c.agent.OnReadError(err)
			continue
		}
		pkgSize := int(binary.LittleEndian.Uint32(buffer))
		if pkgSize == 0 {
			zlog.Debugf("receive ping")
			continue // 心跳包
		}
		readbuf := buffer
		if pkgSize > MaxTcpPkgSize {
			readbuf = make([]byte, pkgSize)
		}
		err = c.Read(readbuf, pkgSize)
		if err != nil {
			c.agent.OnReadError(err)
			continue
		}
		c.agent.OnMessage(readbuf[:pkgSize])
	}
}

func (c *TCPConn) Read(buf []byte, n int) error {
	c.conn.SetReadDeadline(time.Now().Add(MaxTcpWaitTime))
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
