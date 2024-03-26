package agent

import (
	"context"
	"east/core/basic"
	"east/core/util"
	"encoding/binary"
	"io"
	"net"
	"time"
)

const (
	MaxTcpWaitTime = 5 * time.Second
	MaxTcpPkgSize  = 1024
)

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
	basic.Go(func() {
		c.readLoop(ctx)
	})
	basic.Go(func() {

	})
}

func (c *TCPConn) readLoop(ctx context.Context) {
	buffer := make([]byte, MaxTcpPkgSize)
	for {
		if util.ContextDone(ctx) {
			return
		}
		err := c.Read(buffer, 2)
		if err != nil {
			c.agent.OnReadError(err)
			continue
		}
		pkgSize := int(binary.LittleEndian.Uint16(buffer))
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
