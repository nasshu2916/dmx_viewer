//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/domain/repository"
	"github.com/nasshu2916/dmx_viewer/internal/handler/http"
	"github.com/nasshu2916/dmx_viewer/internal/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure"
	"github.com/nasshu2916/dmx_viewer/internal/usecase"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func InitializeTimeHandler(logger *logger.Logger) (*http.TimeHandler, error) {
	wire.Build(
		config.NewConfig,
		infrastructure.NewTimeRepositoryImpl,
		wire.Bind(new(repository.TimeRepository), new(*infrastructure.TimeRepositoryImpl)),
		usecase.NewTimeUseCaseImpl,
		wire.Bind(new(usecase.TimeUseCase), new(*usecase.TimeUseCaseImpl)),
		http.NewTimeHandler,
	)
	return nil, nil
}

func InitializeWebSocketHandler(logger *logger.Logger) (*websocket.WebSocketHandler, error) {
	wire.Build(
		websocket.NewHub,
		websocket.NewWebSocketHandler,
	)
	return nil, nil
}
