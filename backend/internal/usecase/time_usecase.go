package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/beevik/ntp"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/domain/repository"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type TimeUseCase interface {
	GetCurrentTime() time.Time
	StartTimeSync(ctx context.Context)
}

type TimeUseCaseImpl struct {
	logger         *logger.Logger
	timeRepository repository.TimeRepository
	ntpEnabled     bool
	ntpServer      string
	updateInterval time.Duration
	retryCount     int
}

func NewTimeUseCaseImpl(timeRepository repository.TimeRepository, cfg *config.Config, logger *logger.Logger) *TimeUseCaseImpl {
	return &TimeUseCaseImpl{
		logger:         logger,
		timeRepository: timeRepository,
		ntpEnabled:     cfg.NTP.Enabled,
		ntpServer:      cfg.NTP.Server,
		updateInterval: time.Duration(cfg.NTP.UpdateIntervalMinutes) * time.Minute,
		retryCount:     cfg.NTP.RetryCount,
	}
}

func (u *TimeUseCaseImpl) GetCurrentTime() time.Time {
	return u.timeRepository.GetTime()
}

func (u *TimeUseCaseImpl) StartTimeSync(ctx context.Context) {
	if !u.ntpEnabled {
		u.logger.Debug("NTP sync is disabled. Using system time.")
		return
	}

	u.syncTimeFromNTPWithRetry()

	ticker := time.NewTicker(u.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			u.syncTimeFromNTPWithRetry()
		case <-ctx.Done():
			u.logger.Info("Time sync stopped.")
			return
		}
	}
}

func (u *TimeUseCaseImpl) syncTimeFromNTPWithRetry() {
	for i := 0; i <= u.retryCount; i++ {
		response, err := ntp.Query(u.ntpServer)
		if err == nil {
			u.timeRepository.SetQueryResponse(response)
			u.logger.Debug("Time synchronized from NTP server", "server", u.ntpServer, "time", time.Now().Add(response.ClockOffset).String())
			return
		}

		if i < u.retryCount {
			u.logger.Info(
				"Failed to get time from NTP server",
				"server", u.ntpServer,
				"error", err,
				"attempt", fmt.Sprintf("%d/%d", i+1, u.retryCount+1),
			)
			time.Sleep(100 * time.Millisecond)
		}
	}

	u.logger.Error("Failed to synchronize time from NTP server after multiple retries. Using system time.")
}
