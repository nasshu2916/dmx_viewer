package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalHttp "github.com/nasshu2916/dmx_viewer/internal/interface/handler/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTimeUseCase struct {
	mock.Mock
}

func (m *MockTimeUseCase) GetCurrentTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockTimeUseCase) StartTimeSync(ctx context.Context) {
	m.Called(ctx)
}

func TestTimeHandler_GetTime(t *testing.T) {
	mockUseCase := new(MockTimeUseCase)
	expectedTime := time.Date(2024, time.July, 21, 10, 30, 0, 0, time.UTC)
	mockUseCase.On("GetCurrentTime").Return(expectedTime)

	handler := internalHttp.NewTimeHandler(mockUseCase)

	req, err := http.NewRequest("GET", "/time", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	handler.GetTime(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	t.Logf("Response body: %s", rec.Body.String())

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTime.Format(time.RFC3339), response["datetime"])
	mockUseCase.AssertExpectations(t)
}
