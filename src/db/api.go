package db

import (
	"database/sql"
	_ "embed"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func GetItemsToday() []models.GameItem {
	db, err := sql.Open("sqlite3", getDBPath())
	var result = []models.GameItem{}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	currentTime := time.Now()
	rows, err := db.Query("SELECT id, name, prompt, answer FROM prompts ORDER BY id WHERE  playAt=?", currentTime.Format("2006-01-02"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var prompt string
		var answerJson []byte
		err = rows.Scan(&id, &name, &prompt, &answerJson)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, NewItemFromDbEntry(id, name, time.Now(), prompt, answerJson))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetApiData() models.TodayApiResponse {
	var m models.TodayApiResponse = models.TodayApiResponse{Data: map[string]models.GameItem{}}
	data := GetItemsToday()
	for row := range data {
		item := data[row]
		m.Data[item.Id] = item
	}

	return m
}
