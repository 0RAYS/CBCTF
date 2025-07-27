package model

import (
	"github.com/gorilla/websocket"
	"time"
)

type Connection struct {
	Conn       *websocket.Conn
	LastActive time.Time
}
