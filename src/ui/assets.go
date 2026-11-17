package ui

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

var AssetFS embed.FS

func AssetHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Default().Printf("Asset Request path: %s", path)
	if path == "favicon.ico" {
		path = "img/favicon.ico"
	}

	if strings.HasSuffix(path, ".png") {
		w.Header().Add("Content-Type", "image/png")
	} else if strings.HasSuffix(path, ".ico") {
		w.Header().Add("Content-Type", "image/ico")
	} else if strings.HasSuffix(path, ".mjs") {
		w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
	}

	http.ServeFileFS(w, r, AssetFS, path)
}
