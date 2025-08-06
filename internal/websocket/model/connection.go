package model

import (
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn       *websocket.Conn
	LastActive time.Time
}
