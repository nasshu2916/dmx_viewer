package artnet

import (
	"context"
	"time"
)

type ChannelPressureTest struct {
	server *Server
}

func NewChannelPressureTest(server *Server) *ChannelPressureTest {
	return &ChannelPressureTest{server: server}
}

// SimulateHighLoad 高負荷をシミュレート
func (cpt *ChannelPressureTest) SimulateHighLoad(ctx context.Context, packetsPerSecond int, duration time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(packetsPerSecond))
	defer ticker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-ctx.Done():
			cpt.server.logger.Info("Load simulation stopped by context")
			return
		case <-timeout:
			cpt.server.logger.Info("Load simulation completed", "duration", duration)
			return
		case <-ticker.C:
			// ダミーパケットを送信チャンネルに送る
			dummyData := make([]byte, 512)
			if err := cpt.server.SendToWriteChanForTest(dummyData, &DummyAddr{}); err != nil {
				cpt.server.logger.Debug("Failed to send dummy packet during load test", "error", err)
			}
		}
	}
}

func (cpt *ChannelPressureTest) CheckChannelPressure() map[string]interface{} {
	receiveUtil, sendUtil := cpt.server.GetChannelUtilization()
	healthy, msg := cpt.server.IsChannelHealthy()

	return map[string]interface{}{
		"receiveUtilization": receiveUtil,
		"sendUtilization":    sendUtil,
		"isHealthy":          healthy,
		"healthMessage":      msg,
		"droppedReceive":     cpt.server.GetDroppedPackets(),
		"droppedSend":        cpt.server.GetDroppedSendPackets(),
	}
}

// DummyAddr テスト用のダミーアドレス
type DummyAddr struct{}

func (d *DummyAddr) Network() string { return "udp" }
func (d *DummyAddr) String() string  { return "127.0.0.1:6454" }
