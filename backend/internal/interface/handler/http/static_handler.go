package http

import (
	"io/fs"
	"net/http"
)

type StaticHandler struct {
	indexHtml []byte
	assetsFS  fs.FS
}

func NewStaticHandler(indexHtml []byte, assetsFS fs.FS) *StaticHandler {
	return &StaticHandler{
		indexHtml: indexHtml,
		assetsFS:  assetsFS,
	}
}

func (h *StaticHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(h.indexHtml)
}

func (h *StaticHandler) AssetsHandler() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.FS(h.assetsFS)))
}
