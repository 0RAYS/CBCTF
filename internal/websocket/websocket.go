package websocket

import (
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/websocket/handler"
	"CBCTF/internal/websocket/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	AdminClients   = make(map[uint]*model.Connection)
	AdminClientsMu sync.RWMutex
	UserClients    = make(map[uint]*model.Connection)
	UserClientsMu  sync.RWMutex
)

func Init(router *gin.Engine) {
	handler.AddHandler(model.HeartbeatType, handler.KeepAliveHandler)

	router.GET("/ws", middleware.WSAuth, WS)
}

func WS(ctx *gin.Context) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		conn    *websocket.Conn
		msg     []byte
		ok      bool
		err     error
		id      = middleware.GetSelfID(ctx)
	)
	if middleware.IsAdmin(ctx) {
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		mu = &UserClientsMu
		clients = &UserClients
	}
	mu.Lock()
	connection, ok := (*clients)[id]
	mu.Unlock()
	if !ok {
		conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Logger.Warningf("Upgrade error: %s", err)
			return
		}
		defer func(conn *websocket.Conn) {
			if err = conn.Close(); err != nil {
				log.Logger.Warningf("Upgrade error: %s", err)
			}
		}(conn)
		mu.Lock()
		(*clients)[id] = &model.Connection{Conn: conn, LastActive: time.Now()}
		mu.Unlock()
		connection = (*clients)[id]
	}
	for {
		_, msg, err = connection.Conn.ReadMessage()
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Logger.Warningf("Failed to read ws msg: %s", err)
			break
		}
		if err = handler.HandleMessage(connection, msg); err != nil {
			log.Logger.Warningf("Failed to handle ws msg: %s", err)
		}
	}
	mu.Lock()
	delete(*clients, id)
	mu.Unlock()
}
