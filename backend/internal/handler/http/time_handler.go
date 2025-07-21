package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/usecase"
)

type TimeHandler struct {
	timeUseCase usecase.TimeUseCase
}

func NewTimeHandler(timeUseCase usecase.TimeUseCase) *TimeHandler {
	return &TimeHandler{
		timeUseCase: timeUseCase,
	}
}

// StartTimeSync starts NTP time synchronization
func (h *TimeHandler) StartTimeSync(ctx context.Context) {
	go h.timeUseCase.StartTimeSync(ctx)
}

func (h *TimeHandler) GetTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"datetime": h.timeUseCase.GetCurrentTime().Format(time.RFC3339),
	})
}
