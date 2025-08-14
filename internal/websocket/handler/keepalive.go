package handler

import (
	"CBCTF/internal/websocket/model"
	"encoding/json"
	"time"
)

func KeepAliveHandler(conn *model.Connection, msg []byte) error {
	var recv model.HeartbeatRecv
	if err := json.Unmarshal(msg, &recv); err != nil {
		return err
	}
	conn.LastActive = time.Now()
	response := model.HeartbeatSend{
		Type:  model.HeartbeatWSType,
		Msg:   "pong",
		Title: "Heartbeat",
	}
	return conn.Conn.WriteJSON(response)
}
