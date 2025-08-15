package usecase

import (
	"context"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
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

			artPacket, err := packet.Unmarshal(receivedData.Data)
			if err != nil {
				uc.logger.Info("Failed to unmarshal ArtNet packet", "error", err)
				continue
			}

			packet := model.ReceivedArtPacket{
				Packet: artPacket,
				Addr:   receivedData.Addr,
			}

			// パケットを非同期でハンドラーに渡して処理
			uc.packetHandler.HandlePacketAsync(ctx, packet)
		}
	}
}
