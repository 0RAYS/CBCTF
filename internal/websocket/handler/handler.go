package handler

import (
	"CBCTF/internal/websocket/model"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	receiveHandlerMap   = make(map[string]func(*model.Connection, []byte) error)
	receiveHandlerMapMu sync.RWMutex
)

func AddReceiveHandler(requestType string, handler func(*model.Connection, []byte) error) {
	receiveHandlerMapMu.Lock()
	receiveHandlerMap[requestType] = handler
	receiveHandlerMapMu.Unlock()
}

func DeleteReceiveHandler(requestType string) {
	receiveHandlerMapMu.Lock()
	delete(receiveHandlerMap, requestType)
	receiveHandlerMapMu.Unlock()
}

func HandleReceive(conn *model.Connection, msg []byte) error {
	var req model.Receive
	if err := json.Unmarshal(msg, &req); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	handler, exists := receiveHandlerMap[req.Type]
	if !exists {
		return fmt.Errorf("unknown request type: %s", req.Type)
	}

	if err := handler(conn, msg); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}
	return nil
}
