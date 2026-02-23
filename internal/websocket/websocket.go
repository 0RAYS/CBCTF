package websocket

import (
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/websocket/handler"
	"CBCTF/internal/websocket/middleware"
	"CBCTF/internal/websocket/model"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		// 由 middleware.CORS 检查
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	UserClients   = make(map[uint]*model.Connection)
	UserClientsMu sync.RWMutex
)

func Init() {
	handler.AddReceiveHandler(model.HeartbeatWSType, handler.KeepAliveHandler)
}

func WS(ctx *gin.Context) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		id      = middleware.GetSelf(ctx).ID
		ip      = ctx.ClientIP()
	)

	mu = &UserClientsMu
	clients = &UserClients
	// 从请求的 Sec-WebSocket-Protocol 中取第一个协议名回写，
	// 浏览器要求服务端必须选择一个已提议的子协议，否则连接会被关闭。
	var selectedProtocol string
	if proto := ctx.Request.Header.Get("Sec-Websocket-Protocol"); proto != "" {
		parts := strings.SplitN(proto, ",", 2)
		selectedProtocol = strings.TrimSpace(parts[0])
	}
	var responseHeader http.Header
	if selectedProtocol != "" {
		responseHeader = http.Header{"Sec-Websocket-Protocol": {selectedProtocol}}
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, responseHeader)
	if err != nil {
		return
	}

	log.Logger.Infof("New connection from %d %s", id, ip)

	defer func() {
		recover()
		mu.Lock()
		if c, ok := (*clients)[id]; ok {
			_ = c.Conn.Close()
			delete(*clients, id)
		}
		mu.Unlock()
	}()

	c := &model.Connection{Conn: conn, LastActive: time.Now()}

	mu.Lock()
	(*clients)[id] = c
	mu.Unlock()
	prometheus.UpdateWebSocketMetrics(len(*clients))

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		if len(msg) > 0 {
			_ = handler.HandleReceive(c, msg)
		}
	}
	mu.Lock()
	delete(*clients, id)
	mu.Unlock()
	prometheus.UpdateWebSocketMetrics(len(*clients))
}

func Send(id uint, level, t, title, msg string) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
	)
	mu = &UserClientsMu
	clients = &UserClients

	mu.RLock()
	connection, ok := (*clients)[id]
	mu.RUnlock()
	if !ok {
		return
	}
	if err := connection.Conn.WriteJSON(model.Send{Level: level, Type: t, Msg: msg, Title: title}); err == nil {
		log.Logger.Debugf("Send message %s to %d", title, id)
	}
}

func SendToClients(level, t, title, msg string, idL ...uint) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
	)
	mu = &UserClientsMu
	clients = &UserClients

	mu.RLock()
	defer mu.RUnlock()
	var count int
	for _, id := range idL {
		connection, ok := (*clients)[id]
		if !ok {
			continue
		}
		if err := connection.Conn.WriteJSON(model.Send{Level: level, Type: t, Msg: msg, Title: title}); err != nil {
			continue
		}
		count++
	}
	if count > 0 {
		log.Logger.Debugf("Send message %s to %d clients", title, count)
	}
}
