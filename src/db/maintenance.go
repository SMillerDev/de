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
	stmt, err := Db.Prepare("INSERT OR IGNORE INTO prompts(id, name, playAt, prompt, wiki, tags, answer) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	seed := returnSeed()
	for key, val := range seed.Data {
		fmt.Printf("ID: %s, Prompt: %s\n", key, val.Prompts["en"])

		jsonAnswerString, _ := json.Marshal(val.Answer)
		jsonPromptString, _ := json.Marshal(val.Prompts)
		jsonWikiString, _ := json.Marshal(val.Wiki)
		jsonTagString, _ := json.Marshal(val.Tags)
		_, err = stmt.Exec(key, val.Answer.Properties.Name, "1970-01-01", string(jsonPromptString), string(jsonWikiString), string(jsonTagString), string(jsonAnswerString))
		if err != nil {
			log.Fatal(err)
		}
	}
}
