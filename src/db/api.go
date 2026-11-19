package db

import (
	_ "embed"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func GetItemsToday() []models.GameItem {
	var result = []models.GameItem{}
	currentTime := time.Now()
	rows, err := Db.Query("SELECT id, name, prompt, answer FROM prompts WHERE playAt=? ORDER BY id ", currentTime.Format(time.DateOnly))
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

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, currentTime.Format(time.DateOnly), prompt, answerJson))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetApiData() models.TodayApiResponse {
	var m models.TodayApiResponse = models.TodayApiResponse{
		Data:     map[string]models.GameItem{},
		PlayDate: time.Now().Format(time.DateOnly),
	}
	data := GetItemsToday()
	for row := range data {
		item := data[row]
		m.Data[item.Id] = item
	}

	return m
}
