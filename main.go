package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"seanmolenaar.eu/de/src"
)

//go:embed css js img
var assetFS embed.FS

//go:embed html/*
var htmlFS embed.FS

func main() {
	src.InitDB()
	defer src.DbClose()

	src.SetAssetFS(assetFS)
	src.SetHtmlFS(htmlFS)
	http.HandleFunc("/api/{version}/{api}", src.ApiHandler)

	http.HandleFunc("GET /css/", src.AssetHandler)
	http.HandleFunc("GET /js/", src.AssetHandler)
	http.HandleFunc("GET /img/", src.AssetHandler)
	http.HandleFunc("GET /favicon.ico", src.AssetHandler)

	http.HandleFunc("/admin", src.AdminHtmlHandler)
	http.HandleFunc("/admin/{page}", src.AdminHtmlHandler)

	http.HandleFunc("GET /health", src.HealthHandler)

	http.HandleFunc("GET /{$}", src.HtmlHandler)

	port, found := os.LookupEnv("PORT")
	if found == false {
		port = "8085"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}
