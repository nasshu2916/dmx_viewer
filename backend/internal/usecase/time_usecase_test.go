package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/beevik/ntp"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTimeRepository is a mock of TimeRepository interface
type MockTimeRepository struct {
	mock.Mock
}

func (m *MockTimeRepository) GetTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockTimeRepository) ExistsNTPResponse() bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *MockTimeRepository) SetQueryResponse(resp *ntp.Response) {
	m.Called(resp)
}

func TestGetCurrentTime(t *testing.T) {
	mockRepo := new(MockTimeRepository)
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Set up the mock
	mockRepo.On("GetTime").Return(testTime)

	cfg := &config.Config{}           // Not used in GetCurrentTime
	log := *logger.NewLogger("fatal") // Use NewLogger with a level that suppresses output
	uc := NewTimeUseCaseImpl(mockRepo, cfg, log)

	// Execute the method
	resultTime := uc.GetCurrentTime()

	// Assert the result
	assert.Equal(t, testTime, resultTime)
	mockRepo.AssertExpectations(t)
}

func TestStartTimeSync_NTPDisabled(t *testing.T) {
	mockRepo := new(MockTimeRepository)
	cfg := &config.Config{
		NTP: config.NTP{ // Corrected from config.NTPConfig
			Enabled: false,
		},
	}
	log := *logger.NewLogger("fatal") // Use NewLogger with a level that suppresses output
	uc := NewTimeUseCaseImpl(mockRepo, cfg, log)

	// StartTimeSync should return immediately if NTP is disabled
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	uc.StartTimeSync(ctx)

	// Assert that SetQueryResponse was not called
	mockRepo.AssertNotCalled(t, "SetQueryResponse", mock.Anything)
}
