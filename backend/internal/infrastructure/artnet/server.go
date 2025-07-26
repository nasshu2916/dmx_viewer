package artnet

import (
	"fmt"
	"net"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

// SendPacket 送信するパケットの情報
type SendPacket struct {
	Data []byte
	Addr net.Addr
}

type Server struct {
	conn               net.PacketConn
	logger             *logger.Logger
	config             *config.ArtNet
	ipAddress          string
	port               int
	done               chan bool
	receivedChan       chan model.ReceivedPacket // 受信したArtNetパケットを送信するチャネル
	sendChan           chan SendPacket           // 送信するパケットのチャネル
	channelBufferSize  int                       // チャネルのバッファサイズ
	droppedPackets     int64                     // ドロップされたパケット数
	droppedSendPackets int64                     // ドロップされた送信パケット数
}

func NewServer(logger *logger.Logger, cfg *config.ArtNet) *Server {
	channelBufferSize := cfg.ChannelBufferSize
	if channelBufferSize <= 0 {
		channelBufferSize = DefaultChannelBufferSize
	}

	return &Server{
		conn:               nil,
		logger:             logger,
		config:             cfg,
		ipAddress:          "",
		port:               DefaultPort,
		done:               make(chan bool),
		channelBufferSize:  channelBufferSize,
		receivedChan:       make(chan model.ReceivedPacket, channelBufferSize),
		sendChan:           make(chan SendPacket, channelBufferSize),
		droppedPackets:     0,
		droppedSendPackets: 0,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.ipAddress, s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address %s: %w", addr, err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("ArtNet server startup failed: %w", err)
	}
	s.conn = conn

	s.logger.Info("ArtNet server started", "address", addr, "channelBufferSize", s.channelBufferSize)
	pollInterval := time.Duration(s.config.PollIntervalSeconds) * time.Second

	// ArtPollパケットを定期送信するゴルーチンを開始
	pollTicker := time.NewTicker(pollInterval)
	go s.runPollSender(pollTicker)

	// 送信処理を行うゴルーチンを開始
	go s.runSender()

	// 受信処理を行うゴルーチンを開始
	go s.runReceiver()

	// 統計監視を行うゴルーチンを開始（1分間隔）
	statsTicker := time.NewTicker(60 * time.Second)
	go s.runStatMonitor(statsTicker)

	defer func() {
		pollTicker.Stop()
		statsTicker.Stop()

		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
			s.logger.Info("ArtNet server connection closed")
		}
		close(s.receivedChan)
		close(s.sendChan)
	}()

	<-s.done
	return nil
}

func (s *Server) SendToWriteChan(data []byte, addr net.Addr) error {
	if err := s.validateConnection(); err != nil {
		return err
	}

	return s.sendToWriteChanInternal(data, addr)
}

func (s *Server) sendToWriteChanInternal(data []byte, addr net.Addr) error {
	sendPacket := SendPacket{Data: data, Addr: addr}

	select {
	case s.sendChan <- sendPacket:
		return nil
	default:
		queueLength := len(s.sendChan)
		DropPacketWithLog(s.logger, &s.droppedSendPackets, SendChannel, queueLength, s.channelBufferSize, addr.String())
		return fmt.Errorf("%w: %s", ErrSendChannelFull, addr.String())
	}
}

func (s *Server) SendToWriteChanWithTimeout(data []byte, addr net.Addr, timeout time.Duration) error {
	if err := s.validateConnection(); err != nil {
		return err
	}

	return s.sendToWriteChanWithTimeoutInternal(data, addr, timeout)
}

func (s *Server) sendToWriteChanWithTimeoutInternal(data []byte, addr net.Addr, timeout time.Duration) error {
	sendPacket := SendPacket{Data: data, Addr: addr}

	select {
	case s.sendChan <- sendPacket:
		return nil
	case <-time.After(timeout):
		s.logger.Warn("Send channel write timeout",
			"address", addr.String(),
			"timeout", timeout,
			"queueLength", len(s.sendChan),
			"bufferSize", s.channelBufferSize)
		return fmt.Errorf("%w after %v: %s", ErrSendChannelTimeout, timeout, addr.String())
	}
}

func (s *Server) validateConnection() error {
	if s.conn == nil {
		return ErrConnectionNotEstablished
	}
	return nil
}

func (s *Server) Stop() {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
}

func (s *Server) ReceivedChan() <-chan model.ReceivedPacket {
	return s.receivedChan
}

func (s *Server) SendChan() <-chan SendPacket {
	return s.sendChan
}

// テスト専用のチャネル送信メソッド（接続チェックなし）
func (s *Server) SendToWriteChanForTest(data []byte, addr net.Addr) error {
	return s.sendToWriteChanInternal(data, addr)
}

// テスト専用のタイムアウト付きチャネル送信メソッド（接続チェックなし）
func (s *Server) SendToWriteChanWithTimeoutForTest(data []byte, addr net.Addr, timeout time.Duration) error {
	return s.sendToWriteChanWithTimeoutInternal(data, addr, timeout)
}
