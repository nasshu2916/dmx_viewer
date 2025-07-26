package usecase

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/nasshu2916/dmx_viewer/internal/config"
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
	HandlePacket(packet model.ReceivedArtPacket) error
	// ArtNetパケットを非同期で処理する
	HandlePacketAsync(ctx context.Context, packet model.ReceivedArtPacket)
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
	config            *config.ArtNet
	activeGoroutines  int32         // アクティブなゴルーチン数（atomic操作用）
	maxGoroutines     int32         // 最大ゴルーチン数
	processingTimeout time.Duration // 処理タイムアウト
}

// NewArtNetPacketHandler ArtNetPacketHandlerの新しいインスタンスを作成
func NewArtNetPacketHandler(wsUseCase WebSocketUseCase, artNetWriter ArtNetWriter, cfg *config.ArtNet, logger *logger.Logger) *ArtNetPacketHandlerImpl {
	return &ArtNetPacketHandlerImpl{
		wsUseCase:         wsUseCase,
		artNetWriter:      artNetWriter,
		logger:            logger,
		config:            cfg,
		maxGoroutines:     100,             // デフォルト最大ゴルーチン数
		processingTimeout: 5 * time.Second, // デフォルト処理タイムアウト
	}
}

func (h *ArtNetPacketHandlerImpl) HandlePacket(artNetPacket model.ReceivedArtPacket) error {
	switch p := artNetPacket.Packet.(type) {
	case *packet.ArtDMXPacket:
		return h.broadcastDMXPacket(p)
	case *packet.ArtPollPacket:
		return h.handleArtPollPacket(p)
	default:
		h.logger.Debug("Unsupported ArtNet packet type for WebSocket broadcast", "type", artNetPacket.Packet.GetOpCode().String())
		return nil
	}
}

func (h *ArtNetPacketHandlerImpl) broadcastDMXPacket(dmxPacket *packet.ArtDMXPacket) error {
	dmxData, err := model.NewDMXData(dmxPacket)
	if err != nil {
		h.logger.Error("Failed to create DMX data", "error", err)
		return err
	}
	msg := model.NewWebSocketMessage("artnet_dmx_packet", dmxData)
	return h.wsUseCase.BroadcastToTopic("artnet/dmx_packet", msg)
}

// handleArtPollPacket ArtPollパケットを処理し、ArtPollReplyパケットを送信する
func (h *ArtNetPacketHandlerImpl) handleArtPollPacket(_ *packet.ArtPollPacket) error {
	// ArtPollReplyパケットを作成
	replyPacket, err := h.createArtPollReplyPacket()
	if err != nil {
		h.logger.Error("Failed to create ArtPollReply packet", "error", err)
		return err
	}

	// ブロードキャストでArtPollReplyパケットを送信
	if err := h.BroadcastPacket(replyPacket); err != nil {
		h.logger.Error("Failed to broadcast ArtPollReply packet", "error", err)
		return err
	}
	return nil
}

// createArtPollReplyPacket ArtPollReplyパケットを作成する
func (h *ArtNetPacketHandlerImpl) createArtPollReplyPacket() (*packet.ArtPollReplyPacket, error) {
	replyPacket := packet.NewArtPollReplyPacket()

	// IPアドレスを取得
	localIP, err := h.getLocalIPAddress()
	if err != nil {
		h.logger.Warn("Failed to get local IP address, using 127.0.0.1", "error", err)
		localIP = net.ParseIP("127.0.0.1")
	}

	// IPアドレスを設定
	copy(replyPacket.IPAddress[:], localIP.To4())

	// ポート番号を設定
	replyPacket.Port = 6454 // ArtNetの標準ポート

	// VersionInfo (ファームウェアバージョン)
	replyPacket.VersionInfo = 1

	// ショートネームとロングネームを設定
	copy(replyPacket.ShortName[:], []byte(h.config.ShortName))
	copy(replyPacket.LongName[:], []byte(h.config.LongName))

	// ノードのタイプを設定 (Node)
	replyPacket.Style = code.StNode

	// ステータスを設定 (デフォルト値)
	replyPacket.Status1 = code.Status1(0x00)
	replyPacket.Status2 = code.Status2(0x00)

	// ESTAマニュファクチャーコード（適当な値を設定）
	replyPacket.ESTAmanufacturer = [2]byte{'D', 'V'} // DMX Viewer

	// OEMコード
	replyPacket.Oem = 0x0000

	// ポート数（この実装では出力ポートなしで設定）
	replyPacket.NumPorts = 0

	// ノードレポート（NodeReportCodeのスライスとして設定）
	nodeReportStr := "DMX Viewer Ready"
	for i, char := range []byte(nodeReportStr) {
		if i >= len(replyPacket.NodeReport) {
			break
		}
		replyPacket.NodeReport[i] = code.NodeReportCode(char)
	}

	return replyPacket, nil
}

// getLocalIPAddress ローカルIPアドレスを取得する
func (h *ArtNetPacketHandlerImpl) getLocalIPAddress() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// HandlePacketAsync ArtNetパケットを非同期で処理する
func (h *ArtNetPacketHandlerImpl) HandlePacketAsync(ctx context.Context, receivedPacket model.ReceivedArtPacket) {
	// ゴルーチン数の制限をチェック
	currentGoroutines := atomic.LoadInt32(&h.activeGoroutines)

	if currentGoroutines >= h.maxGoroutines {
		h.logger.Warn("Max goroutines reached, dropping packet",
			"activeGoroutines", currentGoroutines,
			"maxGoroutines", h.maxGoroutines,
			"packetType", receivedPacket.Packet.GetOpCode().String())
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
			done <- h.HandlePacket(receivedPacket)
		}()

		select {
		case <-processingCtx.Done():
			h.logger.Warn("Packet processing timed out",
				"timeout", h.processingTimeout,
				"packetType", receivedPacket.Packet.GetOpCode().String())
		case err := <-done:
			if err != nil {
				h.logger.Error("Failed to process packet asynchronously",
					"error", err,
					"packetType", receivedPacket.Packet.GetOpCode().String())
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
