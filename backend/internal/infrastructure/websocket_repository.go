package infrastructure

import (
	"github.com/nasshu2916/dmx_viewer/internal/interface/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// WebSocketRepositoryImpl WebSocketRepositoryの実装
type WebSocketRepositoryImpl struct {
	hub    *websocket.Hub
	logger *logger.Logger
}

func NewWebSocketRepositoryImpl(hub *websocket.Hub, logger *logger.Logger) *WebSocketRepositoryImpl {
	return &WebSocketRepositoryImpl{
		hub:    hub,
		logger: logger,
	}
}

func (r *WebSocketRepositoryImpl) BroadcastToTopic(topic string, message []byte) error {
	r.hub.BroadcastMessage(websocket.SubscribeTopic(topic), message)
	return nil
}

func (r *WebSocketRepositoryImpl) BroadcastToAll(message []byte) error {
	r.hub.BroadcastMessage(websocket.AllSubscribedTopic, message)
	return nil
}
