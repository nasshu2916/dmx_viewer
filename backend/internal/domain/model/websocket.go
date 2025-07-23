package model

import "time"

type WebSocketMessage struct {
	Type      string      `json:"Type"`
	Data      interface{} `json:"Data"`
	Timestamp int64       `json:"Timestamp"`
}

func NewWebSocketMessage(messageType string, data interface{}) *WebSocketMessage {
	return &WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}
