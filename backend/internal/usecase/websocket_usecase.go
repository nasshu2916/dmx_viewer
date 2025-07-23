package usecase

import (
	"encoding/json"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/internal/domain/repository"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// WebSocketUseCase WebSocketに関連するビジネスロジックを定義するインターフェース
type WebSocketUseCase interface {
	// ArtNetパケットをWebSocketクライアントにブロードキャストする
	BroadcastArtNetPacket(packet packet.ArtNetPacket) error
	// 特定のトピックにメッセージをブロードキャストする
	BroadcastToTopic(topic string, data interface{}) error
}

// WebSocketUseCaseImpl WebSocketUseCaseの実装
type WebSocketUseCaseImpl struct {
	wsRepo repository.WebSocketRepository
	logger *logger.Logger
}

// NewWebSocketUseCaseImpl WebSocketUseCaseの新しいインスタンスを作成
func NewWebSocketUseCaseImpl(wsRepo repository.WebSocketRepository, logger *logger.Logger) *WebSocketUseCaseImpl {
	return &WebSocketUseCaseImpl{
		wsRepo: wsRepo,
		logger: logger,
	}
}

// BroadcastArtNetPacket ArtNetパケットをWebSocketクライアントにブロードキャストする
func (uc *WebSocketUseCaseImpl) BroadcastArtNetPacket(artNetPacket packet.ArtNetPacket) error {
	// パケットタイプに応じて処理を分岐
	switch p := artNetPacket.(type) {
	case *packet.ArtDMXPacket:
		return uc.broadcastDMXPacket(p)
	// case *packet.ArtPollReplyPacket:
	// 	return uc.broadcastPollReplyPacket(p)
	default:
		uc.logger.Debug("Unsupported ArtNet packet type", "type", artNetPacket.GetOpCode().String())
		return nil
	}
}

// BroadcastToTopic 特定のトピックにメッセージをブロードキャストする
func (uc *WebSocketUseCaseImpl) BroadcastToTopic(topic string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		uc.logger.Error("Failed to marshal message", "error", err, "topic", topic)
		return err
	}

	if err := uc.wsRepo.BroadcastToTopic(topic, jsonData); err != nil {
		uc.logger.Error("Failed to broadcast message", "error", err, "topic", topic)
		return err
	}
	return nil
}

// broadcastDMXPacket DMXパケットをWebSocketにブロードキャストする
func (uc *WebSocketUseCaseImpl) broadcastDMXPacket(dmxPacket *packet.ArtDMXPacket) error {
	msg := model.NewWebSocketMessage("artnet_dmx_packet", dmxPacket)
	return uc.BroadcastToTopic("artnet/dmx_packet", msg)
}
