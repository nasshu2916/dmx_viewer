package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
	logger   *logger.Logger
}

func NewWebSocketHandler(logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: オリジンを適切にチェックする
				return true
			},
		},
		logger: logger,
	}
}

func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	h.logger.Info("WebSocket connection established.")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("Failed to read WebSocket message: %v", err)
			break
		}
		h.logger.Info("Received message: %s", message)

		// TODO: 受信したメッセージに応じた処理を実装する
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.logger.Error("Failed to write WebSocket message: %v", err)
			break
		}
	}
}
