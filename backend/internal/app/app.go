package app

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/di"
	"github.com/nasshu2916/dmx_viewer/internal/infrastructure/artnet"
	httpHandler "github.com/nasshu2916/dmx_viewer/internal/interface/handler/http"
	"github.com/nasshu2916/dmx_viewer/internal/interface/router"
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

	wsHandler, err := di.InitializeWebSocketHandler(logger)
	if err != nil {
		logger.Fatal("Failed to initialize WebSocket handler: ", err)
	}

	assetsSubFS, err := fs.Sub(assetsFS, "embed_static/assets")
	if err != nil {
		logger.Fatal("Failed to create sub filesystem: ", err)
	}

	artnetServer := artnet.NewServer(logger, &config.ArtNet)

	go timeHandler.StartTimeSync(ctx)
	go wsHandler.Run()
	go func() {
		if err := artnetServer.Run(); err != nil {
			logger.Error("ArtNet server stopped with error: ", err)
		}
	}()

	staticHandler := httpHandler.NewStaticHandler(indexHtml, assetsSubFS)
	router := router.NewRouter(staticHandler, timeHandler, wsHandler)

	logger.Info(fmt.Sprintf("Server started on :%s", config.App.Port))
	http.ListenAndServe(fmt.Sprintf(":%s", config.App.Port), router)
}
