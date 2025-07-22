package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	logger *logger.Logger
	send   chan []byte
	topics map[SubscribeTopic]struct{}
}

type WebSocketMessage struct {
	Type    string          `json:"type"`    // Type of message
	Topic   SubscribeTopic  `json:"topic"`   // Topic name
	Payload json.RawMessage `json:"payload"` // Actual message payload for "publish" type
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

func NewClient(hub *Hub, conn *websocket.Conn, logger *logger.Logger) *Client {
	topics := make(map[SubscribeTopic]struct{})
	topics["default"] = struct{}{}

	return &Client{
		hub:    hub,
		conn:   conn,
		logger: logger,
		send:   make(chan []byte, 256),
		topics: topics,
	}
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("panic recovered in readPump", "panic", r)
		}
		c.hub.LeaveClient(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket read message error", "error", err)
			}
			break
		}

		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			c.logger.Debug("Failed to unmarshal WebSocket message", "error", err, "message", string(message))
			continue
		}

		switch wsMsg.Type {
		case "subscribe":
			c.SubscribeToTopic(wsMsg.Topic)
		case "unsubscribe":
			c.UnsubscribeFromTopic(wsMsg.Topic)
		default:
			c.logger.Debug("Unknown WebSocket message type", "addr", c.conn.RemoteAddr(), "message", wsMsg)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SubscribeToTopic(topic SubscribeTopic) {
	c.hub.subscribe <- SubscribeRequest{
		client: c,
		topic:  topic,
	}
	c.topics[topic] = struct{}{}
}

func (c *Client) UnsubscribeFromTopic(topic SubscribeTopic) {
	c.hub.unsubscribe <- SubscribeRequest{
		client: c,
		topic:  topic,
	}
	delete(c.topics, topic)
}
