package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/zlog"
)

type WebSocketListener struct {
	server   *http.Server
	upgrader websocket.Upgrader
	manager  *AgentManager
}

func GetWebSocketBindAddress() string {
	defaultAddr := fmt.Sprintf(":%d/", conf.ServerID+0x1FEE)
	return conf.String("agent.websocket.address", defaultAddr)
}

func NewWebSocketListener(am *AgentManager) IListener {
	ws := &WebSocketListener{
		manager: am,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	ws.server = &http.Server{
		Handler: ws,
	}
	return ws
}

func (ws *WebSocketListener) Run() error {
	ws.server.Addr = GetWebSocketBindAddress()
	errChan := make(chan error, 1)
	core.Go(func() {
		err := ws.server.ListenAndServe()
		errChan <- err
	})
	time.Sleep(time.Second) // 等待1秒检测端口监听
	zlog.Infof("http listen and serve %s", ws.server.Addr)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (ws *WebSocketListener) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ws.server.Shutdown(ctx)
}

func (ws *WebSocketListener) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	conn, err := ws.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		zlog.Warnf("websocket upgrade error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("websocket upgrade faild"))
		return
	}
	wsconn := &WebSocketConn{
		conn: conn,
	}
	ws.manager.OnConnect(wsconn)
}

type WebSocketConn struct {
	conn *websocket.Conn
}

func (c *WebSocketConn) Write(data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := c.conn.WriteMessage(websocket.BinaryMessage, data)
	return err
}

func (c *WebSocketConn) Read(timeout time.Duration) ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *WebSocketConn) Close() {
	c.conn.Close()
}
