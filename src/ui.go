package src

import (
	"net/http"

	"seanmolenaar.eu/de/src/api"
	"seanmolenaar.eu/de/src/ui"
)

func AssetHandler(w http.ResponseWriter, r *http.Request) {
	ui.AssetHandler(w, r)
}

func AdminHtmlHandler(w http.ResponseWriter, r *http.Request) {
	ui.AdminHtmlHandler(w, r)
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	ui.HtmlHandler(w, r)
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	api.ApiHandler(w, r)
}
