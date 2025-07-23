package artnet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

const DefaultPort = 6454

type Server struct {
	conn       net.PacketConn
	logger     *logger.Logger
	config     *config.ArtNet
	ipAddress  string
	port       int
	done       chan bool
	packetChan chan packet.ArtNetPacket // 受信したArtNetパケットを送信するチャネル
}

func NewServer(logger *logger.Logger, cfg *config.ArtNet) *Server {
	return &Server{
		conn:       nil,
		logger:     logger,
		config:     cfg,
		ipAddress:  "",
		port:       DefaultPort,
		done:       make(chan bool),
		packetChan: make(chan packet.ArtNetPacket), // チャネルを初期化
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

	s.logger.Info("ArtNet server started", "address", addr)
	defer func() {
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
		default:
			n, _, err := s.conn.ReadFrom(buffer)
			if err != nil {
				s.logger.Error("Error reading from ArtNet", "error", err)
				continue
			}

			p, err := packet.Unmarshal(buffer[:n])
			if err != nil {
				s.logger.Info("Failed to unmarshal ArtNet packet", "error", err)
				continue
			}

			s.packetChan <- p
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
