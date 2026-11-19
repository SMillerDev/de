package db

import (
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func returnSeed() models.TodayApiResponse {
	seed, status := os.LookupEnv("DB_SEED")
	if status == false {
		return models.TodayApiResponse{}
	}
	var m models.TodayApiResponse
	json.Unmarshal([]byte(seed), &m)

	return m
}

func WipeDB() {
	os.Remove(getDBPath())
}

func SeedDB() {
	log.Printf("Seeding database")
	stmt, err := Db.Prepare("INSERT OR IGNORE INTO prompts(id, playAt, name, prompt, answer) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	seed := returnSeed()
	for key, val := range seed.Data {
		fmt.Printf("ID: %s, Prompt: %s\n", key, val.Prompt)

		jsonString, _ := json.Marshal(val.Answer)
		_, err = stmt.Exec(key, nil, val.Answer.Properties.Name, val.Prompt, string(jsonString))
		if err != nil {
			log.Fatal(err)
		}
	}
}
