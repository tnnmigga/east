package agent

import (
	"context"
	"east/core/basic"
	"east/core/log"
	"east/core/util"
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
		log.Fatalf("tcpagent listen error %v", err)
	}
	tl := &TCPListener{
		manager:  manager,
		listener: listener,
	}
	basic.Go(func(ctx context.Context) {
		timer := time.NewTicker(MaxTcpWaitTime)
		for {
			select {
			case <-timer.C:
				tl.killDeadAgent()
			case <-ctx.Done():
				return
			}
		}
	})
	log.Infof("tcp listen %s", addr)
	return tl
}

type TCPListener struct {
	manager  *AgentManager
	listener net.Listener
}

func (tl *TCPListener) Run() {
	basic.Go(tl.accept)
}

func (tl *TCPListener) Close() {
	tl.listener.Close()
	for _, agent := range tl.manager.agents {
		agent.conn.Close()
	}
}

func (tl *TCPListener) killDeadAgent() {
	tl.manager.rw.RLock()
	nowNs := util.NowNs()
	for uid, agent := range tl.manager.agents {
		state := atomic.LoadInt32(&agent.state)
		if state == AgentStateDead {
			delete(tl.manager.agents, uid)
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
		delete(tl.manager.agents, uid)
	}
	tl.manager.rw.RUnlock()
}

func (tl *TCPListener) accept() {
	defer util.RecoverPanic()
	for {
		conn, err := tl.listener.Accept()
		if err != nil {
			log.Warnf("tcp accept error %v", err)
			return
		}
		log.Debug("new conn: ", conn.RemoteAddr())
		tcpConn := &TCPConn{
			conn: conn,
		}
		tl.manager.OnConnect(tcpConn)
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
	basic.Go(c.readLoop)
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
			log.Debugf("receive ping")
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
