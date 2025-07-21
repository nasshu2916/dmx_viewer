package websocket

import (
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
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

func NewClient(hub *Hub, conn *websocket.Conn, logger *logger.Logger) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		logger: logger,
		send:   make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
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
		c.logger.Info("Received message from client", "message", message)
		c.hub.BroadcastMessage("default", message)
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
}

func (c *Client) UnsubscribeFromTopic(topic SubscribeTopic) {
	c.hub.unsubscribe <- SubscribeRequest{
		client: c,
		topic:  topic,
	}
}

func (c *Client) GetSubscribedTopics() []SubscribeTopic {
	return c.hub.GetSubscribedTopics(c)
}
