package app

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/internal/di"
	"github.com/nasshu2916/dmx_viewer/internal/handler/websocket"
	"github.com/nasshu2916/dmx_viewer/pkg/httpserver"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

//go:embed "embed_static/index.html"
var indexHtml []byte

//go:embed "embed_static/assets/*"
var assetsFS embed.FS

func Run(ctx context.Context, config *config.Config, logger logger.Logger) {
	var err error

	timeHandler, err := di.InitializeTimeHandler(logger)
	if err != nil {
		logger.Fatal("Failed to initialize time handler: ", err)
	}

	// NTP時刻同期を開始
	go timeHandler.TimeUseCase.StartTimeSync(ctx)

	router := chi.NewRouter()
	wsHandler := websocket.NewWebSocketHandler(&logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHtml)
	})
	var assetsSubFS fs.FS
	if assetsSubFS, err = fs.Sub(assetsFS, "embed_static/assets"); err != nil {
		logger.Fatal("Failed to create sub filesystem: ", err)
	}
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsSubFS))))

	router.Handle("/ws", http.HandlerFunc(wsHandler.ServeWS))

	router.Get("/api/time", http.HandlerFunc(timeHandler.GetTime))

	logger.Info(fmt.Sprintf("Server started on :%s", config.App.Port))
	server := httpserver.New(router, httpserver.Port(config.App.Port))
	if err = <-server.Notify(); err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
