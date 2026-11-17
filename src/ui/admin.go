package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"seanmolenaar.eu/de/src/db"
	"seanmolenaar.eu/de/src/models"
)

func getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
func getAdminUser() string {
	return os.Getenv("ADMIN_USER")
}

func handlePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		log.Printf("Posted for %s with ID %s", r.FormValue("name"), r.FormValue("id"))
		http.Redirect(w, r, "/admin", http.StatusFound)
		return true
	}

	return false
}
func checkedAuth(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/admin") == false {
		return false
	}
	providedUser, providedPassword, ok := r.BasicAuth()

	// Check if the header exists and matches expected credentials
	if ok && providedUser == getAdminUser() && providedPassword == getAdminPassword() {
		return false
	}
	// Return 401 and set the challenge header so browsers prompt for login
	w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
	return true
}

func getItemByID(id string) []models.GameItem {
	var data []models.GameItem = nil

	var err error
	data, err = db.GetItemById(id)
	if err != nil {
		log.Default().Printf("Error when fetching item: %f", err)
		return nil
	}
	return data
}

func AdminHtmlHandler(w http.ResponseWriter, r *http.Request) {
	page := r.PathValue("page")
	log.Default().Printf("Admin Request page: \"%s\" on path: \"%s\"", page, r.URL.Path)
	if checkedAuth(w, r) {
		return
	}

	if handlePost(w, r) {
		return
	}

	var data []models.GameItem = nil
	switch page {
	case "":
		page = "html/admin/index.html"
		data = db.GetDBItems(10000)
	default:
		page = fmt.Sprintf("html/admin/%s.html", page)
	}
	id := r.URL.Query().Get("id")
	if id != "" {
		data = getItemByID(id)
		if len(data) == 0 {
			log.Default().Printf("Did not find item: %s", id)
			http.NotFound(w, r)
			return
		}
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
		log.Default().Printf("Error when parsing page: %f", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
