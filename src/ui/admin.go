package ui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"uuid"

	"seanmolenaar.eu/de/src/db"
	"seanmolenaar.eu/de/src/models"
)

func getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
func getAdminUser() string {
	return os.Getenv("ADMIN_USER")
}

func getLanguages() []string {
	langString, found := os.LookupEnv("LANGUAGES")
	if found == false {
		langString = "en"
	}

	return strings.Split(langString, ",")
}

func handleDelete(w http.ResponseWriter, r *http.Request) bool {
	page := r.PathValue("page")
	if page != "delete" {
		return false
	}
	var id string
	if id = r.FormValue("id"); id == "" {
		http.Error(w, "No ID provided!", http.StatusBadRequest)
		return true
	}

	db.DeleteItemById(id)
	http.Redirect(w, r, "/admin", http.StatusFound)
	return true
}
func handlePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	var id string
	if id = r.FormValue("id"); id == "" {
		id = uuid.New().String()
	}

	log.Printf("Posted for %s with ID %s", r.FormValue("name"), id)
	var answer models.GeoJSONAnswer
	err := json.Unmarshal([]byte(r.FormValue("answer")), &answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return true
	}
	answer.Properties.Name = r.FormValue("name")
	prompts := map[string]string{}
	wikis := map[string]string{}
	languages := getLanguages()
	for langIndex := range languages {
		prompts[languages[langIndex]] = r.FormValue(fmt.Sprintf("prompt-%s", languages[langIndex]))
		wikis[languages[langIndex]] = r.FormValue(fmt.Sprintf("wiki-%s", languages[langIndex]))
	}
	item := models.GameItem{
		Id:      id,
		PlayAt:  time.Time{},
		Prompts: prompts,
		Answer:  answer,
		Tags:    strings.Split(r.FormValue("tags"), ","),
		Wiki:    wikis,
	}
	if dateString := r.FormValue("playat"); dateString != "" {
		date, err := time.Parse(time.DateOnly, dateString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return true
		}
		item.PlayAt = date
	}
	err = db.SaveItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
	return true
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

func AdminHtmlHandler(w http.ResponseWriter, r *http.Request) {
	page := r.PathValue("page")
	log.Default().Printf("Admin Request page: \"%s\" on path: \"%s\"", page, r.URL.Path)
	if checkedAuth(w, r) {
		return
	}

	if handlePost(w, r) {
		return
	}

	if handleDelete(w, r) {
		return
	}

	data := AdminPageData{Languages: getLanguages(), Now: time.Now()}
	where := make(map[string]string)
	id := r.FormValue("id")
	if r.FormValue("tag") != "" {
		where["tags"] = r.FormValue("tag")
		log.Printf("Filtering on tag: %s", where["tags"])
	} else if r.FormValue("date") != "" {
		where["playAt"] = r.FormValue("date")
		log.Printf("Filtering on date: %s", where["playAt"])
	} else if id != "" {
		where["id"] = r.FormValue("id")
		log.Printf("Filtering on id: %s", id)
	}
	data.Items = db.GetDBItems(10000, where)
	switch page {
	case "":
		page = "html/admin/index.html"
		data.Dates = db.GetDateCounts()
		data.Tags = db.GetTagCounts()
	case "dump.json":
		page = fmt.Sprintf("html/admin/%s.html", page)
	default:
		page = fmt.Sprintf("html/admin/%s.html", page)
	}
	if len(data.Items) == 0 && id != "" {
		log.Default().Printf("Did not find item: %s", id)
		http.NotFound(w, r)
		return
	}
	if r.PathValue("page") == "reports" {
		data.Reports = db.GetReports(id, 100)
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	t, err := template.New(path.Base(page)).Funcs(template.FuncMap{
		"lower":       strings.ToLower,
		"StringsJoin": strings.Join,
		"json": func(s any) string {
			bytes, err := json.Marshal(s)
			if err != nil {
				return "FAILED TO MARSHAL JSON!"
			}
			return string(bytes)
		},
	}).ParseFS(HtmlFS, page)
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

type AdminPageData struct {
	Languages []string
	Now       time.Time
	Items     []models.GameItem
	Reports   []models.ReportValues
	Tags      map[string]int
	Dates     map[string]int
}
