package usecase

import (
	"encoding/json"

	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/internal/domain/repository"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// WebSocketUseCase WebSocketに関連するビジネスロジックを定義するインターフェース
type WebSocketUseCase interface {
	BroadcastToTopic(topic string, message *model.WebSocketMessage) error
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

// BroadcastToTopic 特定のトピックにメッセージをブロードキャストする
func (uc *WebSocketUseCaseImpl) BroadcastToTopic(topic string, message *model.WebSocketMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		uc.logger.Error("Failed to marshal WebSocket message", "error", err, "topic", topic)
		return err
	}

	if err := uc.wsRepo.BroadcastToTopic(topic, jsonData); err != nil {
		uc.logger.Error("Failed to broadcast WebSocket message", "error", err, "topic", topic)
		return err
	}
	return nil
}
