package ui

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

var HtmlFS embed.FS

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Default().Printf("HTML Request path: %s", path)

	page := "html/index.html"
	var data MainPageData = MainPageData{
		Language: getSelectedLanguage(r),
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFS(HtmlFS, page)
	if err != nil {
		log.Default().Printf("Error when parsing page: %f", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error!"))
		return
	}
	if t == nil {
		log.Default().Printf("Did not find page: %s", page)
		http.NotFound(w, r)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

type MainPageData struct {
	Language string
}

func getSelectedLanguage(r *http.Request) string {
	if r.Header.Get("Accept-Language") != "" {
		log.Printf("Accepting languages: %s", r.Header.Get("Accept-Language"))
	}
	return getLanguages()[0]
}
