package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/usecase"
)

type TimeHandler struct {
	TimeUseCase usecase.TimeUseCase
}

func NewTimeHandler(timeUseCase usecase.TimeUseCase) *TimeHandler {
	return &TimeHandler{
		TimeUseCase: timeUseCase,
	}
}

func (h *TimeHandler) GetTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"datetime": h.TimeUseCase.GetCurrentTime().Format(time.RFC3339),
	})
}
