package main

import (
	"embed"
	"log"
	"net/http"

	"seanmolenaar.eu/de/src"
)

//go:embed css js img
var assetFS embed.FS

//go:embed html/*
var htmlFS embed.FS

func main() {
	src.WipeDB()
	src.SeedDB()
	src.SetAssetFS(assetFS)
	src.SetHtmlFS(htmlFS)
	http.HandleFunc("/api/", src.ApiHandler)

	http.HandleFunc("GET /css/", src.AssetHandler)
	http.HandleFunc("GET /js/", src.AssetHandler)
	http.HandleFunc("GET /img/", src.AssetHandler)
	http.HandleFunc("GET /favicon.ico", src.AssetHandler)

	http.HandleFunc("/admin", src.AdminHtmlHandler)
	http.HandleFunc("/admin/{page}", src.AdminHtmlHandler)

	http.HandleFunc("GET /{$}", src.HtmlHandler)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
