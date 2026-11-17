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
	path := r.URL.Path[1:]
	log.Default().Printf("API Request path: %s", path)
	switch path {
	case "api/v1/result":
		report(r)
		w.WriteHeader(http.StatusNoContent)
		return
	case "api/v1/base.geo.json":
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(baseData())
		return
	case "api/v1/today.json":
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
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 32)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 32)
	var report = models.ReportValues{
		Id:        r.URL.Query().Get("id"),
		Lattitude: lat,
		Longitude: lng,
	}
	log.Default().Printf("Reported values: id=%s", report.Id)
}
