package websocket

import (
	"encoding/json"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

const AllSubscribedTopic = "BROADCAST_ALL"

var dmxIncrement = 0

type SubscribeTopic string

type TopicMessage struct {
	topic   SubscribeTopic
	message []byte
}

type SubscribeRequest struct {
	topic  SubscribeTopic
	client *Client
}

type Hub struct {
	logger *logger.Logger

	clients           map[*Client]struct{}                    // All registered clients
	SubscribedClients map[SubscribeTopic]map[*Client]struct{} // Registered subscribed clients, grouped by topic.

	join        chan *Client          // Channel for new client joining the hub
	leave       chan *Client          // Channel for client leaving the hub
	subscribe   chan SubscribeRequest // Channel for subscribing to topics
	unsubscribe chan SubscribeRequest // Channel for unsubscribing from topics

	broadcast chan TopicMessage // Channel for broadcasting messages to subscribed clients
}

func NewHub(logger *logger.Logger) *Hub {
	return &Hub{
		logger: logger,

		clients:           make(map[*Client]struct{}),
		SubscribedClients: make(map[SubscribeTopic]map[*Client]struct{}),

		join:        make(chan *Client),
		leave:       make(chan *Client),
		subscribe:   make(chan SubscribeRequest),
		unsubscribe: make(chan SubscribeRequest),

		broadcast: make(chan TopicMessage),
	}
}

func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic recovered in Hub.Run", "panic", r)
		}
	}()

	testPollingTicker := time.NewTicker(200 * time.Millisecond) // 500ms interval for test polling
	defer testPollingTicker.Stop()
	go func() {
		for range testPollingTicker.C {
			h.testPollingDmxMessage()
		}
	}()

	for {
		select {
		case client := <-h.join:
			h.clients[client] = struct{}{}
			for t := range client.topics {
				h.subscribeTopic(client, t)
			}
			h.logger.Debug("Client joined", "addr", client.conn.RemoteAddr())

		case client := <-h.leave:
			for t := range client.topics {
				h.unsubscribeTopic(client, t)
			}
			delete(h.clients, client)
			close(client.send)
			h.logger.Debug("Client left", "addr", client.conn.RemoteAddr())

		case request := <-h.subscribe:
			h.subscribeTopic(request.client, request.topic)

		case request := <-h.unsubscribe:
			h.unsubscribeTopic(request.client, request.topic)

		case topicMessage := <-h.broadcast:
			if clientsInTopic, ok := h.SubscribedClients[topicMessage.topic]; ok {
				for client := range clientsInTopic {
					select {
					case client.send <- topicMessage.message:
					default:
						if _, ok := h.clients[client]; ok {
							h.logger.Info("Failed to send message to client, closing connection", "addr", client.conn.RemoteAddr(), "topic", topicMessage.topic)
						} else {
							h.logger.Warn("Failed to send message to client, client not found", "addr", client.conn.RemoteAddr(), "topic", topicMessage.topic)
							// If the client is not connected, remove it from the hub
							h.unsubscribeTopic(client, topicMessage.topic)
						}
					}
				}
			}
		}
	}
}

func (h *Hub) subscribeTopic(client *Client, topic SubscribeTopic) {
	if _, ok := h.SubscribedClients[topic]; !ok {
		h.SubscribedClients[topic] = make(map[*Client]struct{})
	}
	h.SubscribedClients[topic][client] = struct{}{}
}

func (h *Hub) unsubscribeTopic(client *Client, topic SubscribeTopic) {
	if clientsInTopic, ok := h.SubscribedClients[topic]; ok {
		delete(clientsInTopic, client)
		if len(clientsInTopic) == 0 {
			delete(h.SubscribedClients, topic)
		}
	}
}

func (h *Hub) JoinClient(client *Client) {
	h.join <- client
}

func (h *Hub) LeaveClient(client *Client) {
	h.leave <- client
}

func (h *Hub) BroadcastMessage(topic SubscribeTopic, message []byte) {
	h.broadcast <- TopicMessage{topic: topic, message: message}
}

func (h *Hub) testPollingDmxMessage() {
	// テスト用のポーリングメッセージを送信
	dmxPacket := packet.NewArtDMXPacket()
	// ランダムな length 512 の DMX データを生成
	var dmxDataArr [512]byte
	for i := range dmxDataArr {
		dmxDataArr[i] = byte((i + dmxIncrement) % 256) // 0-255
	}
	dmxPacket.Data = dmxDataArr
	dmxIncrement++

	dmxData, err := model.NewDMXData(dmxPacket)
	if err != nil {
		h.logger.Error("Failed to create DMX data", "error", err)
		return
	}
	msg := model.NewWebSocketMessage("artnet_dmx_packet", dmxData)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal WebSocket message", "error", err)
		return
	}
	h.BroadcastMessage(AllSubscribedTopic, jsonData)
}
