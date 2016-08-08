package artnet

import (
	"fmt"
	"sync/atomic"
	"time"
)

// runStatMonitor 統計監視を行うゴルーチン
func (s *Server) runStatMonitor(statsTicker *time.Ticker) {
	panicHandler := NewPanicHandler(s.logger, "stat_monitor")
	defer panicHandler.Handle()
	defer statsTicker.Stop()

	for {
		select {
		case <-s.done:
			s.logger.Debug("Stat monitor stopped")
			return
		case <-statsTicker.C:
			s.logChannelStatistics()
		}
	}
}

// logChannelStatistics チャンネル統計をログ出力
func (s *Server) logChannelStatistics() {
	bufferSize, receiveQueueLength, sendQueueLength, droppedReceivePackets, droppedSendPackets := s.GetChannelStats()

	// 受信チャンネルの統計
	receiveUtil := CalculateUtilization(receiveQueueLength, bufferSize)
	receiveStatus := DetermineHealthStatus(receiveUtil, droppedReceivePackets)
	LogChannelStats(s.logger, ReceiveChannel, receiveQueueLength, bufferSize, droppedReceivePackets, receiveStatus)

	// 送信チャンネルの統計
	sendUtil := CalculateUtilization(sendQueueLength, bufferSize)
	sendStatus := DetermineHealthStatus(sendUtil, droppedSendPackets)
	LogChannelStats(s.logger, SendChannel, sendQueueLength, bufferSize, droppedSendPackets, sendStatus)
} // GetChannelStats チャネルの統計情報を取得

func (s *Server) GetChannelStats() (bufferSize int, receiveQueueLength int, sendQueueLength int, droppedReceivePackets int64, droppedSendPackets int64) {
	return s.channelBufferSize, len(s.receivedChan), len(s.sendChan), atomic.LoadInt64(&s.droppedPackets), atomic.LoadInt64(&s.droppedSendPackets)
}

// GetDroppedPackets ドロップされた受信パケット数を取得
func (s *Server) GetDroppedPackets() int64 {
	return atomic.LoadInt64(&s.droppedPackets)
}

// GetDroppedSendPackets ドロップされた送信パケット数を取得
func (s *Server) GetDroppedSendPackets() int64 {
	return atomic.LoadInt64(&s.droppedSendPackets)
}

// ResetDroppedPackets ドロップされたパケット数をリセット
func (s *Server) ResetDroppedPackets() {
	atomic.StoreInt64(&s.droppedPackets, 0)
	atomic.StoreInt64(&s.droppedSendPackets, 0)
}

// IsChannelHealthy チャンネルの健全性をチェック
func (s *Server) IsChannelHealthy() (bool, string) {
	_, receiveQueueLength, sendQueueLength, droppedReceivePackets, droppedSendPackets := s.GetChannelStats()

	receiveUtil := float64(receiveQueueLength) / float64(s.channelBufferSize) * 100
	sendUtil := float64(sendQueueLength) / float64(s.channelBufferSize) * 100

	if receiveUtil > 90 || sendUtil > 90 {
		return false, fmt.Sprintf("Critical utilization: receive=%.1f%%, send=%.1f%%", receiveUtil, sendUtil)
	}

	if droppedReceivePackets > 0 || droppedSendPackets > 0 {
		return false, fmt.Sprintf("Packets dropped: receive=%d, send=%d", droppedReceivePackets, droppedSendPackets)
	}

	if receiveUtil > 75 || sendUtil > 75 {
		return false, fmt.Sprintf("High utilization: receive=%.1f%%, send=%.1f%%", receiveUtil, sendUtil)
	}

	return true, "Channels are healthy"
}

// GetChannelUtilization チャンネル使用率を取得
func (s *Server) GetChannelUtilization() (receiveUtil float64, sendUtil float64) {
	_, receiveQueueLength, sendQueueLength, _, _ := s.GetChannelStats()
	receiveUtil = float64(receiveQueueLength) / float64(s.channelBufferSize) * 100
	sendUtil = float64(sendQueueLength) / float64(s.channelBufferSize) * 100
	return
}

// 受信パケット数（トータル）と直近1分の集計を管理するためのフィールドと関数群

// recordPacketAt は指定時刻の受信を1件記録する
func (s *Server) recordPacketAt(now time.Time) {
	// 総受信数をインクリメント
	atomic.AddInt64(&s.packetsReceivedTotal, 1)

	// 秒単位のリングバッファ（60秒）にインクリメント
	sec := now.Unix()
	idx := sec % 60
	// バケットの時刻が異なる場合はリセットして現在秒に合わせる
	prevSec := atomic.LoadInt64(&s.packetsReceivedBucketSec[idx])
	if prevSec != sec {
		// 他ゴルーチンとの競合を考慮してCASで時刻を更新し、成功した場合のみリセット
		if atomic.CompareAndSwapInt64(&s.packetsReceivedBucketSec[idx], prevSec, sec) {
			atomic.StoreInt64(&s.packetsReceivedBuckets[idx], 0)
		}
	}
	atomic.AddInt64(&s.packetsReceivedBuckets[idx], 1)
}

// recordReceivedPacket は現在時刻で1件の受信を記録する
func (s *Server) recordReceivedPacket() { s.recordPacketAt(time.Now()) }

// GetReceivedPacketsTotal は起動以降の総受信パケット数を返す
func (s *Server) GetReceivedPacketsTotal() int64 {
	return atomic.LoadInt64(&s.packetsReceivedTotal)
}

// GetReceivedPacketsLastMinute は直近60秒間に受信したパケット数を返す
func (s *Server) GetReceivedPacketsLastMinute() int64 {
	nowSec := time.Now().Unix()
	var sum int64 = 0
	for i := int64(0); i < 60; i++ {
		sec := atomic.LoadInt64(&s.packetsReceivedBucketSec[i])
		if nowSec-sec < 60 {
			sum += atomic.LoadInt64(&s.packetsReceivedBuckets[i])
		}
	}
	return sum
}

// RecordPacketAtForTest テスト用ヘルパー（他パッケージのテストから呼べるよう公開）
func (s *Server) RecordPacketAtForTest(t time.Time) { s.recordPacketAt(t) }
