package usecase

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// ArtNetPacketHandler ArtNetパケットを処理するハンドラーのインターフェース
type ArtNetPacketHandler interface {
	// ArtNetパケットを処理する
	HandlePacket(packet packet.ArtNetPacket) error
	// ArtNetパケットを非同期で処理する
	HandlePacketAsync(ctx context.Context, packet packet.ArtNetPacket)
	// 処理中のゴルーチン数を取得
	GetActiveGoroutines() int
}

// ArtNetPacketHandlerImpl ArtNetPacketHandlerの実装
type ArtNetPacketHandlerImpl struct {
	wsUseCase         WebSocketUseCase
	logger            *logger.Logger
	activeGoroutines  int32         // アクティブなゴルーチン数（atomic操作用）
	maxGoroutines     int32         // 最大ゴルーチン数
	processingTimeout time.Duration // 処理タイムアウト
}

// NewArtNetPacketHandler ArtNetPacketHandlerの新しいインスタンスを作成
func NewArtNetPacketHandler(wsUseCase WebSocketUseCase, logger *logger.Logger) *ArtNetPacketHandlerImpl {
	return &ArtNetPacketHandlerImpl{
		wsUseCase:         wsUseCase,
		logger:            logger,
		maxGoroutines:     100,             // デフォルト最大ゴルーチン数
		processingTimeout: 5 * time.Second, // デフォルト処理タイムアウト
	}
}

// NewArtNetPacketHandlerWithConfig 設定付きでArtNetPacketHandlerの新しいインスタンスを作成
func NewArtNetPacketHandlerWithConfig(wsUseCase WebSocketUseCase, logger *logger.Logger, maxGoroutines int32, timeout time.Duration) *ArtNetPacketHandlerImpl {
	return &ArtNetPacketHandlerImpl{
		wsUseCase:         wsUseCase,
		logger:            logger,
		maxGoroutines:     maxGoroutines,
		processingTimeout: timeout,
	}
}

// HandlePacket ArtNetパケットを処理する
func (h *ArtNetPacketHandlerImpl) HandlePacket(artNetPacket packet.ArtNetPacket) error {
	switch p := artNetPacket.(type) {
	case *packet.ArtDMXPacket:
		dmxData, err := model.NewDMXData(p)
		if err != nil {
			h.logger.Error("Failed to create DMX data", "error", err)
			return err
		}
		return h.broadcastDMXPacket(dmxData)
	default:
		h.logger.Debug("Unsupported ArtNet packet type for WebSocket broadcast", "type", artNetPacket.GetOpCode().String())
		return nil
	}
}

func (h *ArtNetPacketHandlerImpl) broadcastDMXPacket(dmxData *model.DMXData) error {
	msg := model.NewWebSocketMessage("artnet_dmx_packet", dmxData)
	return h.wsUseCase.BroadcastToTopic("artnet/dmx_packet", msg)
}

// HandlePacketAsync ArtNetパケットを非同期で処理する
func (h *ArtNetPacketHandlerImpl) HandlePacketAsync(ctx context.Context, artNetPacket packet.ArtNetPacket) {
	// ゴルーチン数の制限をチェック
	currentGoroutines := atomic.LoadInt32(&h.activeGoroutines)

	if currentGoroutines >= h.maxGoroutines {
		h.logger.Warn("Max goroutines reached, dropping packet",
			"activeGoroutines", currentGoroutines,
			"maxGoroutines", h.maxGoroutines,
			"packetType", artNetPacket.GetOpCode().String())
		return
	}

	// ゴルーチンカウンターを増加
	atomic.AddInt32(&h.activeGoroutines, 1)

	// 非同期でパケット処理を実行
	go func() {
		defer func() {
			// ゴルーチンカウンターを減少
			atomic.AddInt32(&h.activeGoroutines, -1)

			if r := recover(); r != nil {
				h.logger.Error("Panic occurred in async packet processing", "panic", r)
			}
		}()

		// タイムアウト付きコンテキストを作成
		processingCtx, cancel := context.WithTimeout(ctx, h.processingTimeout)
		defer cancel()

		// タイムアウト制御付きで処理を実行
		done := make(chan error, 1)
		go func() {
			done <- h.HandlePacket(artNetPacket)
		}()

		select {
		case <-processingCtx.Done():
			h.logger.Warn("Packet processing timed out",
				"timeout", h.processingTimeout,
				"packetType", artNetPacket.GetOpCode().String())
		case err := <-done:
			if err != nil {
				h.logger.Error("Failed to process packet asynchronously",
					"error", err,
					"packetType", artNetPacket.GetOpCode().String())
			}
		}
	}()
}

// GetActiveGoroutines 処理中のゴルーチン数を取得
func (h *ArtNetPacketHandlerImpl) GetActiveGoroutines() int {
	return int(atomic.LoadInt32(&h.activeGoroutines))
}
