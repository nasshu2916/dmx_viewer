package usecase

import (
	"encoding/json"

	"github.com/nasshu2916/dmx_viewer/internal/domain/repository"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// WebSocketUseCase WebSocketに関連するビジネスロジックを定義するインターフェース
type WebSocketUseCase interface {
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
