package usecase

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// ArtNetWriter ArtNetパケットを送信するためのインターフェース
type ArtNetWriter interface {
	SendToWriteChan(data []byte, addr net.Addr) error
}

// ArtNetPacketHandler ArtNetパケットを処理するハンドラーのインターフェース
type ArtNetPacketHandler interface {
	// ArtNetパケットを処理する
	HandlePacket(packet packet.ArtNetPacket) error
	// ArtNetパケットを非同期で処理する
	HandlePacketAsync(ctx context.Context, packet packet.ArtNetPacket)
	// 処理中のゴルーチン数を取得
	GetActiveGoroutines() int
	// ArtNetパケットを送信する
	SendPacket(artNetPacket packet.ArtNetPacket, addr net.Addr) error
	// ArtNetパケットをブロードキャストする
	BroadcastPacket(artNetPacket packet.ArtNetPacket) error
}

// ArtNetPacketHandlerImpl ArtNetPacketHandlerの実装
type ArtNetPacketHandlerImpl struct {
	wsUseCase         WebSocketUseCase
	artNetWriter      ArtNetWriter
	logger            *logger.Logger
	activeGoroutines  int32         // アクティブなゴルーチン数（atomic操作用）
	maxGoroutines     int32         // 最大ゴルーチン数
	processingTimeout time.Duration // 処理タイムアウト
}

// NewArtNetPacketHandler ArtNetPacketHandlerの新しいインスタンスを作成
func NewArtNetPacketHandler(wsUseCase WebSocketUseCase, artNetWriter ArtNetWriter, logger *logger.Logger) *ArtNetPacketHandlerImpl {
	return &ArtNetPacketHandlerImpl{
		wsUseCase:         wsUseCase,
		artNetWriter:      artNetWriter,
		logger:            logger,
		maxGoroutines:     100,             // デフォルト最大ゴルーチン数
		processingTimeout: 5 * time.Second, // デフォルト処理タイムアウト
	}
}

// NewArtNetPacketHandlerWithConfig 設定付きでArtNetPacketHandlerの新しいインスタンスを作成
func NewArtNetPacketHandlerWithConfig(wsUseCase WebSocketUseCase, artNetWriter ArtNetWriter, logger *logger.Logger, maxGoroutines int32, timeout time.Duration) *ArtNetPacketHandlerImpl {
	return &ArtNetPacketHandlerImpl{
		wsUseCase:         wsUseCase,
		artNetWriter:      artNetWriter,
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

// SendPacket ArtNetパケットを指定されたアドレスに送信する
// 使用例:
//
//	dmxPacket := packet.NewArtDMXPacket()
//	dmxPacket.Universe = 0
//	dmxPacket.Data = []byte{255, 128, 0, ...} // DMXデータ
//	addr := &net.UDPAddr{IP: net.ParseIP("192.168.1.100"), Port: 6454}
//	err := handler.SendPacket(dmxPacket, addr)
func (h *ArtNetPacketHandlerImpl) SendPacket(artNetPacket packet.ArtNetPacket, addr net.Addr) error {
	data, err := artNetPacket.MarshalBinary()
	if err != nil {
		h.logger.Error("Failed to marshal ArtNet packet for sending", "error", err, "packetType", artNetPacket.GetOpCode().String())
		return err
	}

	if err := h.artNetWriter.SendToWriteChan(data, addr); err != nil {
		h.logger.Error("Failed to send ArtNet packet", "error", err, "address", addr.String(), "packetType", artNetPacket.GetOpCode().String())
		return err
	}

	h.logger.Debug("Successfully sent ArtNet packet", "address", addr.String(), "packetType", artNetPacket.GetOpCode().String())
	return nil
}

// BroadcastPacket ArtNetパケットをブロードキャストする
// 使用例:
//
//	pollPacket := packet.NewArtPollPacket()
//	err := handler.BroadcastPacket(pollPacket)
func (h *ArtNetPacketHandlerImpl) BroadcastPacket(artNetPacket packet.ArtNetPacket) error {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 6454, // ArtNetのデフォルトポート
	}

	return h.SendPacket(artNetPacket, broadcastAddr)
}
