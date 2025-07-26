package usecase

import (
	"context"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
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
	packetHandler ArtNetPacketHandler
	logger        *logger.Logger
}

// NewArtNetUseCaseImpl ArtNetBridgeUseCaseの新しいインスタンスを作成
func NewArtNetUseCaseImpl(packetHandler ArtNetPacketHandler, logger *logger.Logger) *ArtNetBridgeUseCaseImpl {
	return &ArtNetBridgeUseCaseImpl{
		packetHandler: packetHandler,
		logger:        logger,
	}
}

// StartPacketForwarding ArtNetサーバーからのパケットをWebSocketに転送する処理を開始
func (uc *ArtNetBridgeUseCaseImpl) StartPacketForwarding(ctx context.Context, artNetServer *artnet.Server) {
	defer func() {
		if r := recover(); r != nil {
			uc.logger.Error("Panic occurred in packet forwarding", "panic", r)
		}
	}()

	receivedChan := artNetServer.ReceivedChan()

	uc.logger.Info("Started ArtNet packet forwarding to WebSocket")

	// パフォーマンス監視用のゴルーチンを開始
	go uc.monitorPerformance(ctx)

	for {
		select {
		case <-ctx.Done():
			uc.logger.Info("Packet forwarding stopped due to context cancellation")
			return

		case receivedData, ok := <-receivedChan:
			if !ok {
				uc.logger.Info("ArtNet packet channel closed, stopping packet forwarding")
				return
			}

			packet, err := packet.Unmarshal(receivedData.Data)
			if err != nil {
				uc.logger.Info("Failed to unmarshal ArtNet packet", "error", err)
				continue
			}

			// パケットを非同期でハンドラーに渡して処理
			uc.packetHandler.HandlePacketAsync(ctx, packet)
		}
	}
}

// パフォーマンスを監視する
func (uc *ArtNetBridgeUseCaseImpl) monitorPerformance(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // 30秒間隔でモニタリング
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			activeGoroutines := uc.packetHandler.GetActiveGoroutines()
			if activeGoroutines > 0 {
				uc.logger.Info("ArtNet packet processing performance",
					"activeGoroutines", activeGoroutines)
			}
		}
	}
}
