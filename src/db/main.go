package db

import (
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func getDBDir() string {
	datadir := os.Getenv("DATADIR")
	if datadir != "" {
		return datadir
	}
	return "."
}

func getDBPath() string {
	return fmt.Sprintf("%s/data.db", getDBDir())
}

func NewItemFromDbEntry(id string, name string, playAt time.Time, prompt string, answerBytes []byte) models.GameItem {
	var answer models.GeoJSONAnswer
	json.Unmarshal(answerBytes, &answer)
	answer.Properties.Name = name
	return models.GameItem{Id: id, Prompt: prompt, Answer: answer, PlayAt: playAt}
}
