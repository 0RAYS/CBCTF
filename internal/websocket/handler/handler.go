package handler

import (
	"CBCTF/internal/websocket/model"
	"encoding/json"
	"fmt"
)

var (
	receiveHandlerMap = make(map[string]func(*model.Connection, []byte) error)
)

func AddReceiveHandler(requestType string, handler func(*model.Connection, []byte) error) {
	receiveHandlerMap[requestType] = handler
}

func HandleReceive(conn *model.Connection, msg []byte) error {
	var recv model.Recv
	if err := json.Unmarshal(msg, &recv); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	handler, ok := receiveHandlerMap[recv.Type]
	if !ok {
		return fmt.Errorf("unknown request type: %s", recv.Type)
	}
	if err := handler(conn, msg); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}
	return nil
}
