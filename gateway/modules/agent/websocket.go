package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/web"
	"github.com/tnnmigga/nett/zlog"
)

type WebSocketListener struct {
	addr     string
	ha       *web.HttpAgent
	upgrader websocket.Upgrader
	manager  *AgentManager
}

func GetWebSocketBindAddress() string {
	defaultAddr := fmt.Sprintf(":%d/", conf.ServerID()+0x1FEE)
	return conf.String("agent.tcp.addr", defaultAddr)
}

func NewWebSocketListener(am *AgentManager) IListener {
	addr := GetWebSocketBindAddress()
	ws := &WebSocketListener{
		ha:      web.NewHttpAgent(),
		manager: am,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		addr: addr,
	}
	return ws
}

func (ws *WebSocketListener) Run() {
	ws.ha.GET(ws.addr, ws.handle)
	ws.ha.Run(ws.addr)
}

func (ws *WebSocketListener) handle(ctx *gin.Context) {
	conn, err := ws.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zlog.Warnf("websocket upgrade error %v", err)
		ctx.String(http.StatusInternalServerError, "websocket upgrade faild")
		return
	}
	wsconn := &WebSocketConn{
		conn: conn,
	}
	ws.manager.OnConnect(wsconn)
}

func (ws *WebSocketListener) Close() {
	ws.ha.Stop()
	for _, agent := range ws.manager.agents {
		agent.conn.Close()
	}
}

type WebSocketConn struct {
	agent IAgent
	conn  *websocket.Conn
}

func (c *WebSocketConn) Write(data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := c.conn.WriteMessage(websocket.BinaryMessage, data)
	return err
}

func (c *WebSocketConn) Run(ctx context.Context) {
	core.Go(c.readLoop)
}

func (c *WebSocketConn) readLoop() {
	for {
		c.conn.SetReadDeadline(time.Now().Add(MaxAliveTime))
		mType, data, err := c.conn.ReadMessage()
		if err != nil {
			c.agent.OnReadError(err)
		}
		if mType == websocket.PingMessage {
			c.conn.WriteMessage(websocket.PongMessage, nil)
			continue
		}
		if mType == websocket.BinaryMessage {
			c.agent.OnMessage(data)
			continue
		}
		zlog.Warnf("invalid msg type %v", mType)
	}
}

func (c *WebSocketConn) Close() {
	c.conn.Close()
}

func (c *WebSocketConn) BindAgent(agent IAgent) {
	c.agent = agent
}
