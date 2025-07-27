package handler

import (
	"CBCTF/internal/websocket/model"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	handlerMap   = make(map[string]func(*model.Connection, []byte) error)
	handlerMapMu sync.RWMutex
)

func AddHandler(requestType string, handler func(*model.Connection, []byte) error) {
	handlerMapMu.Lock()
	handlerMap[requestType] = handler
	handlerMapMu.Unlock()
}

func DeleteHandler(requestType string) {
	handlerMapMu.Lock()
	delete(handlerMap, requestType)
	handlerMapMu.Unlock()
}

func HandleMessage(conn *model.Connection, msg []byte) error {
	var req model.Receive
	if err := json.Unmarshal(msg, &req); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	handler, exists := handlerMap[req.Type]
	if !exists {
		return fmt.Errorf("unknown request type: %s", req.Type)
	}

	if err := handler(conn, msg); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}
	return nil
}
