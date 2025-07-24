package usecase

import (
	"context"

	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// ArtNetBridgeUseCase ArtNetサーバーとWebSocketの橋渡しを行うビジネスロジック
type ArtNetBridgeUseCase interface {
	// ArtNetサーバーからのパケットをWebSocketに転送する処理を開始
	StartPacketForwarding(ctx context.Context, artnetServer *artnet.Server) error
}

// ArtNetBridgeUseCaseImpl ArtNetBridgeUseCaseの実装
type ArtNetBridgeUseCaseImpl struct {
	wsUseCase WebSocketUseCase
	logger    *logger.Logger
}

// NewArtNetUseCaseImpl ArtNetBridgeUseCaseの新しいインスタンスを作成
func NewArtNetUseCaseImpl(wsUseCase WebSocketUseCase, logger *logger.Logger) *ArtNetBridgeUseCaseImpl {
	return &ArtNetBridgeUseCaseImpl{
		wsUseCase: wsUseCase,
		logger:    logger,
	}
}

// StartPacketForwarding ArtNetサーバーからのパケットをWebSocketに転送する処理を開始
func (uc *ArtNetBridgeUseCaseImpl) StartPacketForwarding(ctx context.Context, artNetServer *artnet.Server) {
	defer func() {
		if r := recover(); r != nil {
			uc.logger.Error("Panic occurred in packet forwarding", "panic", r)
		}
	}()

	packetChan := artNetServer.PacketChan()

	uc.logger.Info("Started ArtNet packet forwarding to WebSocket")

	for {
		select {
		case <-ctx.Done():
			uc.logger.Info("Packet forwarding stopped due to context cancellation")
			return

		case artnetPacket, ok := <-packetChan:
			if !ok {
				uc.logger.Info("ArtNet packet channel closed, stopping packet forwarding")
				return
			}

			// パケットをWebSocketにブロードキャスト
			if err := uc.wsUseCase.BroadcastArtNetPacket(artnetPacket); err != nil {
				uc.logger.Error("Failed to broadcast ArtNet packet", "error", err, "opcode", artnetPacket.GetOpCode())
			}
		}
	}
}
