package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/interface/httpctx"
	"github.com/nasshu2916/dmx_viewer/internal/usecase"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

type TimeHandler struct {
	timeUseCase usecase.TimeUseCase
	logger      *logger.Logger
}

func NewTimeHandler(timeUseCase usecase.TimeUseCase, logger *logger.Logger) *TimeHandler {
	return &TimeHandler{
		timeUseCase: timeUseCase,
		logger:      logger,
	}
}

// StartTimeSync starts NTP time synchronization
func (h *TimeHandler) StartTimeSync(ctx context.Context) {
	go h.timeUseCase.StartTimeSync(ctx)
}

func (h *TimeHandler) GetTime(w http.ResponseWriter, r *http.Request) {
	// アクセスログ（Request-ID/Real-IP）
	h.logger.Info("time handler: GetTime",
		"request_id", r.Header.Get("X-Request-Id"),
		"real_ip", httpctx.RealIP(r.Context()),
		"method", r.Method,
		"path", r.URL.Path,
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"datetime": h.timeUseCase.GetCurrentTime().Format(time.RFC3339),
	})
}
