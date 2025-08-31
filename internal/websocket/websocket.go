package websocket

import (
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/websocket/handler"
	"CBCTF/internal/websocket/middleware"
	"CBCTF/internal/websocket/model"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	AdminClients   = make(map[uint]*model.Connection)
	AdminClientsMu sync.RWMutex
	UserClients    = make(map[uint]*model.Connection)
	UserClientsMu  sync.RWMutex
)

func Init() {
	handler.AddReceiveHandler(model.HeartbeatWSType, handler.KeepAliveHandler)
}

func WS(ctx *gin.Context) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		role    string
		id      = middleware.GetSelfID(ctx)
		ip      = ctx.ClientIP()
	)

	if middleware.IsAdmin(ctx) {
		role = "admin"
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		role = "user"
		mu = &UserClientsMu
		clients = &UserClients
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger.Warningf("Failed to upgrade connection: %s", err)
		return
	}

	log.Logger.Infof("New connection from %s %d %s", role, id, ip)

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

	c := &model.Connection{Conn: conn, LastActive: time.Now()}

	mu.Lock()
	(*clients)[id] = c
	mu.Unlock()
	prometheus.UpdateWebSocketMetrics(len(*clients))

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Logger.Warningf("Failed to read ws msg: %s", err)
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
	prometheus.UpdateWebSocketMetrics(len(*clients))
}

func Send(admin bool, id uint, level, t, title, msg string) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		role    string
	)
	if admin {
		role = "admin"
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		role = "user"
		mu = &UserClientsMu
		clients = &UserClients
	}

	mu.RLock()
	connection, ok := (*clients)[id]
	mu.RUnlock()
	if !ok {
		log.Logger.Warningf("Failed to found %s-%d connection", role, id)
		return
	}
	if err := connection.Conn.WriteJSON(model.Send{Level: level, Type: t, Msg: msg, Title: title}); err != nil {
		log.Logger.Warningf("Failed to send message %s to %s %d: %s", title, role, id, err)
	} else {
		log.Logger.Debugf("Send message %s to %s %d", title, role, id)
	}
}

func SendToClients(admin bool, level, t, title, msg string, idL ...uint) {
	var (
		mu      *sync.RWMutex
		clients *map[uint]*model.Connection
		role    string
	)
	if admin {
		role = "admin"
		mu = &AdminClientsMu
		clients = &AdminClients
	} else {
		role = "user"
		mu = &UserClientsMu
		clients = &UserClients
	}

	mu.RLock()
	defer mu.RUnlock()
	var count int
	for _, id := range idL {
		connection, ok := (*clients)[id]
		if !ok {
			continue
		}
		if err := connection.Conn.WriteJSON(model.Send{Level: level, Type: t, Msg: msg, Title: title}); err != nil {
			log.Logger.Warningf("Failed to send message %s to %s %d: %s", title, role, id, err)
			continue
		}
		count++
	}
	if count > 0 {
		log.Logger.Debugf("Send message %s to %s %d clients", title, role, count)
	} else {
		log.Logger.Warningf("Failed to send message %s to %s %d clients", title, role, len(idL))
	}
}
