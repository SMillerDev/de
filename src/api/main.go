package api

import (
	_ "embed"
	"encoding/json/v2"
	"log"
	"net/http"
	"os"
	"strconv"

	"seanmolenaar.eu/de/src/db"
	"seanmolenaar.eu/de/src/models"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	if version != "v1" {
		http.Error(w, "Invalid API version!", http.StatusBadRequest)
		return
	}
	api := r.PathValue("api")
	log.Default().Printf("API Request path: %s", r.URL.Path)
	switch api {
	case "result":
		report(r)
		w.WriteHeader(http.StatusNoContent)
		return
	case "base.geo.json":
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(baseData())
		return
	case "today.json":
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		m := db.GetApiData()
		if err := json.MarshalWrite(w, &m); err != nil {
			// Inability to marshal the output suggests a server-side problem.
			// This error is not always observable by the client since
			// json.MarshalWrite may have already written to the output.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.NotFound(w, r)
	}
}

func baseData() []byte {
	base, status := os.LookupEnv("GEO_BASE_JSON")
	if status == false {
		return []byte("{}")
	}
	return []byte(base)
}

func report(r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 32)
	lng, _ := strconv.ParseFloat(r.FormValue("lng"), 32)
	var report = models.ReportValues{
		Id:        r.FormValue("id"),
		Lattitude: lat,
		Longitude: lng,
	}
	log.Default().Printf("Reported values: id=%s", report.Id)
}
