package agent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/tnnmigga/core/conc"
	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/infra/zlog"
)

const (
	pkgHeadSizeByteLen = 4
	defaultBufferLen   = 1024
	defaultBindPort    = 9527
)

func GetTCPBindAddress() string {
	defaultAddr := fmt.Sprintf(":%d", defaultBindPort)
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
	conc.Go(func() {
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
	})
	zlog.Infof("tcp listen %s", addr)
	return nil
}

func (tcp *TCPListener) Close() {
	tcp.listener.Close()
}

type TCPConn struct {
	conn net.Conn
	wBuf []byte
	rBuf []byte
}

func NewTCPConn(conn net.Conn) *TCPConn {
	return &TCPConn{
		conn: conn,
		wBuf: make([]byte, defaultBufferLen),
		rBuf: make([]byte, pkgHeadSizeByteLen), // 只读取包头
	}
}

func (c *TCPConn) Write(data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	buf := c.wBuf
	pkgSize := len(data) + pkgHeadSizeByteLen
	if pkgSize < defaultBufferLen {
		buf = make([]byte, pkgSize)
	}
	binary.LittleEndian.PutUint32(buf, uint32(len(data)))
	copy(buf[pkgHeadSizeByteLen:], data)
	_, err := c.conn.Write(buf[:pkgSize])
	zlog.Debugf("agent write %v", data)
	return err
}

func (c *TCPConn) Read(timeout time.Duration) ([]byte, error) {
	err := c.readFull(c.rBuf, pkgHeadSizeByteLen, timeout)
	if err != nil {
		return nil, err
	}
	pkgSize := int(binary.LittleEndian.Uint16(c.rBuf))
	if pkgSize == 0 {
		zlog.Debugf("receive ping")
		return nil, nil
	}
	// 每次创建一个新缓冲区
	// 防止在传递过程中可能出现的slice并发读写
	buffer := make([]byte, pkgSize)
	err = c.readFull(buffer, pkgSize, timeout)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (c *TCPConn) readFull(buf []byte, n int, timeout time.Duration) error {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	if _, err := io.ReadFull(c.conn, buf[:n]); err != nil {
		return err
	}
	return nil
}

func (c *TCPConn) Close() {
	c.conn.Close()
}
