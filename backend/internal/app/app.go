package app

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/di"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	metrics "github.com/nasshu2916/dmx_viewer/internal/infrastructure/metrics"
	httpHandler "github.com/nasshu2916/dmx_viewer/internal/interface/handler/http"
	"github.com/nasshu2916/dmx_viewer/internal/interface/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/internal/interface/router"
	"github.com/nasshu2916/dmx_viewer/internal/usecase"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

//go:embed "embed_static/index.html"
var indexHtml []byte

//go:embed "embed_static/assets/*"
var assetsFS embed.FS

func Run(ctx context.Context, config *config.Config, logger *logger.Logger) {
	timeHandler, err := di.InitializeTimeHandler(logger)
	if err != nil {
		logger.Fatal("Failed to initialize time handler: ", err)
	}

	hub := websocket.NewHub(logger)
	go hub.Run()

	wsHandler := websocket.NewWebSocketHandler(hub, logger)

	// HubからWebSocketRepositoryとUseCaseを作成
	wsRepo := infrastructure.NewWebSocketRepositoryImpl(hub, logger)
	wsUseCase := usecase.NewWebSocketUseCaseImpl(wsRepo, logger)

	artNetServer := artnet.NewServer(logger, &config.ArtNet)
	artNetNodeRepo := infrastructure.NewArtNetNodeRepository()
	artNetPacketHandler := usecase.NewArtNetPacketHandler(wsUseCase, artNetServer, &config.ArtNet, logger, artNetNodeRepo)
	artNetUseCase := usecase.NewArtNetUseCaseImpl(artNetPacketHandler, logger)

	assetsSubFS, err := fs.Sub(assetsFS, "embed_static/assets")
	if err != nil {
		logger.Fatal("Failed to create sub filesystem: ", err)
	}

	go timeHandler.StartTimeSync(ctx)
	go func() {
		if err := artNetServer.Run(); err != nil {
			logger.Error("ArtNet server stopped with error: ", err)
		}
	}()

	// ArtNetパケットをWebSocketに転送する処理を開始
	go artNetUseCase.StartPacketForwarding(ctx, artNetServer)

	staticHandler := httpHandler.NewStaticHandler(indexHtml, assetsSubFS, logger)
	healthHandler := httpHandler.NewHealthHandler(artNetServer, logger)

	// Prometheus メトリクス登録（プロセス/Go標準 + ArtNet カスタム）
	metrics.RegisterArtNetMetrics(artNetServer, logger)
	metricsHandler := httpHandler.NewMetricsHandler(nil, logger)

	httpTimeout := time.Duration(config.App.HTTPTimeoutSeconds) * time.Second
	router := router.NewRouter(staticHandler, timeHandler, healthHandler, metricsHandler, wsHandler, logger, httpTimeout)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.App.Port),
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server started on :%s", config.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server stopped with error: ", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error: ", err)
	} else {
		logger.Info("HTTP server shutdown gracefully")
	}
}
