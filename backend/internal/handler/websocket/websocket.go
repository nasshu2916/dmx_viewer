package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
	hub      *Hub
	logger   *logger.Logger
}

func NewWebSocketHandler(hub *Hub, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: オリジンを適切にチェックする
				return true
			},
		},
		hub:    hub,
		logger: logger,
	}
}

func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	client := &Client{hub: h.hub, conn: conn, send: make(chan []byte, 256), logger: h.logger}
	client.hub.RegisterClient(client)

	go client.writePump()
	go client.readPump()
}

func (h *WebSocketHandler) Run() {
	go h.hub.Run()
}
