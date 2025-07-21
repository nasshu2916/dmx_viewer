package websocket

import (
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

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

	clients      map[SubscribeTopic]map[*Client]bool // Registered subscribed clients, grouped by topic.
	clientTopics map[*Client]map[SubscribeTopic]bool // Topics subscribed by each client

	broadcast chan TopicMessage

	subscribe   chan SubscribeRequest // Channel for subscribing to topics
	unsubscribe chan SubscribeRequest // Channel for unsubscribing from topics

	unregister chan *Client // Channel for unregistering clients
}

func NewHub(logger *logger.Logger) *Hub {
	return &Hub{
		logger:       logger,
		clients:      make(map[SubscribeTopic]map[*Client]bool),
		clientTopics: make(map[*Client]map[SubscribeTopic]bool),
		broadcast:    make(chan TopicMessage),
		subscribe:    make(chan SubscribeRequest),
		unsubscribe:  make(chan SubscribeRequest),
		unregister:   make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case request := <-h.subscribe:
			h.subscribeTopic(request.client, request.topic)

		case request := <-h.unsubscribe:
			h.unsubscribeTopic(request.client, []SubscribeTopic{request.topic})

		case client := <-h.unregister:
			h.unregisterAndUnsubscribeClient(client)

		case topicMessage := <-h.broadcast:
			if clientsInTopic, ok := h.clients[topicMessage.topic]; ok {
				for client := range clientsInTopic {
					select {
					case client.send <- topicMessage.message:
					default:
						h.logger.Info("Broadcast topic failed", "addr", client.conn.RemoteAddr(), "topic", topicMessage.topic)
					}
				}
			}
		}
	}
}

func (h *Hub) subscribeTopic(client *Client, topic SubscribeTopic) {
	if _, ok := h.clients[topic]; !ok {
		h.clients[topic] = make(map[*Client]bool)
	}
	h.clients[topic][client] = true

	// Update clientTopics
	if _, ok := h.clientTopics[client]; !ok {
		h.clientTopics[client] = make(map[SubscribeTopic]bool)
	}
	h.clientTopics[client][topic] = true
}

func (h *Hub) unsubscribeTopic(client *Client, topics []SubscribeTopic) {
	for _, topic := range topics {
		if clientsInTopic, ok := h.clients[topic]; ok {
			delete(clientsInTopic, client)
			if len(clientsInTopic) == 0 {
				delete(h.clients, topic)
			}
		}
		// Update clientTopics for each topic
		if clientTopics, ok := h.clientTopics[client]; ok {
			delete(clientTopics, topic)
			if len(clientTopics) == 0 {
				delete(h.clientTopics, client)
			}
		}
	}
}

func (h *Hub) unregisterAndUnsubscribeClient(client *Client) {
	topics := h.GetSubscribedTopics(client)
	h.unsubscribeTopic(client, topics)
	close(client.send)
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastMessage(topic SubscribeTopic, message []byte) {
	h.broadcast <- TopicMessage{topic: topic, message: message}
}

func (h *Hub) GetSubscribedTopics(client *Client) []SubscribeTopic {
	var result []SubscribeTopic
	if topics, ok := h.clientTopics[client]; ok {
		for topic := range topics {
			result = append(result, topic)
		}
	}
	return result
}
