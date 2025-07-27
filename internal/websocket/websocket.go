package websocket

import (
	"CBCTF/internal/config"
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
			origin := r.Header.Get("Origin")
			if config.Env.Frontend == "*" {
				return true
			}
			return origin == config.Env.Frontend
		},
	}

	AdminClients   = make(map[uint]*model.Connection)
	AdminClientsMu sync.RWMutex
	UserClients    = make(map[uint]*model.Connection)
	UserClientsMu  sync.RWMutex
)

func Init() {
	handler.AddReceiveHandler(model.HeartbeatType, handler.KeepAliveHandler)
}

func WS(ctx *gin.Context) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		id      = middleware.GetSelfID(ctx)
	)

	if middleware.IsAdmin(ctx) {
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		mu = &UserClientsMu
		clients = &UserClients
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger.Warningf("Upgrade error: %s", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Errorf("Recovered in WS handler: %v", r)
		}
		mu.Lock()
		if c, ok := (*clients)[id]; ok {
			_ = c.Conn.Close()
			delete(*clients, id)
		}
		mu.Unlock()
	}()

	mu.Lock()
	(*clients)[id] = &model.Connection{Conn: conn, LastActive: time.Now()}
	mu.Unlock()

	for {
		mu.RLock()
		c := (*clients)[id]
		mu.RUnlock()

		_, msg, err := c.Conn.ReadMessage()
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Logger.Warningf("Failed to read ws msg: %v", err)
			break
		}
		if len(msg) > 0 {
			if err := handler.HandleReceive(c, msg); err != nil {
				log.Logger.Debugf("Failed to handle ws msg %s: %s", msg, err)
			}
		}
	}
	mu.Lock()
	delete(*clients, id)
	mu.Unlock()
}

func Send(admin bool, id uint, level, t, title, msg string) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
	)
	if admin {
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		mu = &UserClientsMu
		clients = &UserClients
	}

	mu.RLock()
	connection, ok := (*clients)[id]
	mu.RUnlock()
	if !ok {
		log.Logger.Warningf("No connection found with ID %d", id)
		return
	}
	sendMsg := model.Send{
		Level: level,
		Type:  t,
		Msg:   msg,
		Title: title,
	}
	if err := connection.Conn.WriteJSON(sendMsg); err != nil {
		log.Logger.Warningf("Failed to send message ID %d: %s", id, err)
	}
}

func SendToAll(admin bool, level, t, title, msg string) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
	)
	if admin {
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		mu = &UserClientsMu
		clients = &UserClients
	}

	mu.RLock()
	defer mu.RUnlock()
	for id, connection := range *clients {
		sendMsg := model.Send{
			Level: level,
			Type:  t,
			Msg:   msg,
			Title: title,
		}
		if err := connection.Conn.WriteJSON(sendMsg); err != nil {
			log.Logger.Warningf("Failed to send message to ID %d: %s", id, err)
		}
	}
}
