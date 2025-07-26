package artnet

import (
	"fmt"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
)

// runReceiver 受信処理を行うゴルーチン
func (s *Server) runReceiver() {
	panicHandler := NewPanicHandler(s.logger, "receiver")
	defer panicHandler.Handle()

	buffer := make([]byte, DefaultMaxPacketSize)

	for {
		select {
		case <-s.done:
			s.logger.Debug("Receiver stopped")
			return
		default:
			if err := s.processIncomingPackets(buffer); err != nil {
				// エラーが発生してもログに記録して受信を続ける
				s.logger.Debug("Processing incoming packet error, continuing", "error", err)
			}
		}
	}
}

// processIncomingPackets 受信パケットを処理
func (s *Server) processIncomingPackets(buffer []byte) error {
	s.conn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	n, recievedAddr, err := s.conn.ReadFrom(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// タイムアウトの場合、done チャンネルをチェックしてから続行
			select {
			case <-s.done:
				s.logger.Debug("Receiver stopped during timeout")
				return fmt.Errorf("receiver stopped")
			default:
				// タイムアウトは正常な動作なので、デバッグレベルでログ出力して続行
				s.logger.Debug("Read timeout, continuing to receive")
				return nil // エラーを返さずに受信を続ける
			}
		}
		// ネットワークエラーもログに記録するが受信は続ける
		s.logger.Warn("Error reading from ArtNet, continuing to receive", "error", err)
		return nil // エラーを返さずに受信を続ける
	}

	data := make([]byte, n)
	copy(data, buffer[:n])

	receivedPacket := model.ReceivedPacket{
		Data: data,
		Addr: recievedAddr,
	}

	return s.sendToReceiveChannel(receivedPacket)
}

// sendToReceiveChannel 受信チャンネルにパケットを送信
func (s *Server) sendToReceiveChannel(packet model.ReceivedPacket) error {
	select {
	case s.receivedChan <- packet:
		return nil
	default:
		queueLength := len(s.receivedChan)
		DropPacketWithLog(s.logger, &s.droppedPackets, ReceiveChannel, queueLength, s.channelBufferSize, packet.Addr.String())
		// チャンネルが満杯でもパケットを破棄して受信を続ける
		return nil
	}
}

// runSender 送信チャネルからパケットを受信して送信するゴルーチン
func (s *Server) runSender() {
	panicHandler := NewPanicHandler(s.logger, "sender")
	defer panicHandler.Handle()

	for {
		select {
		case <-s.done:
			s.logger.Debug("Sender stopped")
			return
		case sendPacket := <-s.sendChan:
			s.processSendPacket(sendPacket)
		}
	}
}

// processSendPacket 送信パケットを処理
func (s *Server) processSendPacket(sendPacket SendPacket) {
	if s.conn == nil {
		s.logger.Error("Connection is not established, dropping packet", "address", sendPacket.Addr.String())
		return
	}

	n, err := s.conn.WriteTo(sendPacket.Data, sendPacket.Addr)
	if err != nil {
		s.logger.Error("Error writing to ArtNet", "error", err, "address", sendPacket.Addr.String())
		return
	}

	s.logger.Debug("Sent packet", "to", sendPacket.Addr.String(), "size", n)

	// 送信チャンネルの使用率をチェック
	s.checkSendChannelUtilization()
}

// checkSendChannelUtilization 送信チャンネルの使用率をチェック
func (s *Server) checkSendChannelUtilization() {
	queueLength := len(s.sendChan)
	utilization := CalculateUtilization(queueLength, s.channelBufferSize)

	if utilization > HighUtilizationThreshold {
		status := DetermineHealthStatus(utilization, 0)
		LogChannelStats(s.logger, SendChannel, queueLength, s.channelBufferSize, 0, status)
	}
}

// runPollSender ArtPollパケットを定期的に送信するゴルーチン
func (s *Server) runPollSender(pollTicker *time.Ticker) {
	defer pollTicker.Stop()

	pollPacket := packet.NewArtPollPacket()
	data, err := pollPacket.MarshalBinary()
	if err != nil {
		s.logger.Error("Failed to marshal ArtPoll packet", "error", err)
		return
	}

	broadcastAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: s.port}

	for {
		select {
		case <-s.done:
			s.logger.Debug("Poll sender stopped")
			return
		case <-pollTicker.C:
			s.sendArtPollPacket(data, broadcastAddr)
		}
	}
}

// sendArtPollPacket ArtPollパケットを送信
func (s *Server) sendArtPollPacket(data []byte, addr *net.UDPAddr) {
	if err := s.SendToWriteChan(data, addr); err != nil {
		s.logger.Warn("Failed to send ArtPoll packet", "error", err)
	}
}
