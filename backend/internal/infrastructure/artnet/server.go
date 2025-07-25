package artnet

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

const DefaultPort = 6454
const DefaultChannelBufferSize = 1000

type Server struct {
	conn              net.PacketConn
	logger            *logger.Logger
	config            *config.ArtNet
	ipAddress         string
	port              int
	done              chan bool
	packetChan        chan packet.ArtNetPacket // 受信したArtNetパケットを送信するチャネル
	channelBufferSize int                      // チャネルのバッファサイズ
	droppedPackets    int64                    // ドロップされたパケット数
}

func NewServer(logger *logger.Logger, cfg *config.ArtNet) *Server {
	channelBufferSize := cfg.ChannelBufferSize
	if channelBufferSize <= 0 {
		channelBufferSize = DefaultChannelBufferSize
	}

	return &Server{
		conn:              nil,
		logger:            logger,
		config:            cfg,
		ipAddress:         "",
		port:              DefaultPort,
		done:              make(chan bool),
		channelBufferSize: channelBufferSize,
		packetChan:        make(chan packet.ArtNetPacket, channelBufferSize),
		droppedPackets:    0,
	}
}

func (s *Server) Run() error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Panic occurred in ArtNet server", "panic", r)
		}
	}()

	addr := fmt.Sprintf("%s:%d", s.ipAddress, s.port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("ArtNet server startup failed: %w", err)
	}
	s.conn = conn

	s.logger.Info("ArtNet server started", "address", addr, "channelBufferSize", s.channelBufferSize)
	pollInterval := time.Duration(s.config.PollIntervalSeconds) * time.Second
	pollTicker := time.NewTicker(pollInterval)

	// 統計情報出力用のタイマー（1分間隔）
	statsTicker := time.NewTicker(60 * time.Second)

	defer func() {
		pollTicker.Stop()
		statsTicker.Stop()

		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
			s.logger.Info("ArtNet server connection closed")
		}
		close(s.packetChan)
	}()

	buffer := make([]byte, 1500)
	for {
		select {
		case <-s.done:
			return nil
		case <-statsTicker.C:
			// 統計情報をログ出力
			bufferSize, queueLength, droppedPackets := s.GetChannelStats()
			if droppedPackets > 0 || queueLength > bufferSize/2 {
				s.logger.Info("ArtNet server statistics",
					"bufferSize", bufferSize,
					"queueLength", queueLength,
					"droppedPackets", droppedPackets,
					"utilization", float64(queueLength)/float64(bufferSize)*100)
			}
		case <-pollTicker.C:
			pollPacket := packet.NewArtPollPacket()
			data, err := pollPacket.MarshalBinary()
			if err != nil {
				s.logger.Error("Failed to marshal ArtPoll packet", "error", err)
				continue
			}

			// ブロードキャストアドレスに送信
			broadcastAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: s.port}
			_, err = s.Write(data, broadcastAddr)
			if err != nil {
				s.logger.Error("Failed to send ArtPoll packet", "error", err)
			}
		default:
			// 受信処理はノンブロッキングで行う
			s.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, _, err := s.conn.ReadFrom(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // タイムアウトの場合は次のループへ
				}
				s.logger.Error("Error reading from ArtNet", "error", err)
				continue
			}

			p, err := packet.Unmarshal(buffer[:n])
			if err != nil {
				s.logger.Info("Failed to unmarshal ArtNet packet", "error", err)
				continue
			}

			select {
			case s.packetChan <- p:
			default:
				// チャネルが満杯の場合、パケットをドロップ
				dropped := atomic.AddInt64(&s.droppedPackets, 1)
				if dropped%100 == 0 { // 100パケットごとにログ出力
					s.logger.Warn("ArtNet packet channel is full, dropping packets",
						"droppedPackets", dropped,
						"channelBufferSize", s.channelBufferSize,
						"packetType", p.GetOpCode().String())
				}
			}
		}
	}
}

func (s *Server) Write(data []byte, addr net.Addr) (int, error) {
	if s.conn == nil {
		return 0, fmt.Errorf("ArtNet connection is not established")
	}

	n, err := s.conn.WriteTo(data, addr)
	if err != nil {
		s.logger.Error("Error writing to ArtNet", "error", err, "address", addr.String())
		return 0, err
	}
	s.logger.Debug("Sent packet", "to", addr.String(), "size", n)
	return n, nil
}

func (s *Server) Stop() {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
}

func (s *Server) PacketChan() <-chan packet.ArtNetPacket {
	return s.packetChan
}

// ドロップされたパケット数を取得
func (s *Server) GetDroppedPackets() int64 {
	return atomic.LoadInt64(&s.droppedPackets)
}

// チャネルの統計情報を取得
func (s *Server) GetChannelStats() (bufferSize int, queueLength int, droppedPackets int64) {
	return s.channelBufferSize, len(s.packetChan), atomic.LoadInt64(&s.droppedPackets)
}

// ドロップされたパケット数をリセット
func (s *Server) ResetDroppedPackets() {
	atomic.StoreInt64(&s.droppedPackets, 0)
}
