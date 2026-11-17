package db

import (
	"database/sql"
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
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS prompts (id PRIMARY KEY, playAt, name NOT NULL, prompt NOT NULL, answer json UNIQUE NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO prompts(id, name, prompt, answer) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	seed := returnSeed()
	for key, val := range seed.Data {
		fmt.Printf("ID: %s, Prompt: %s\n", key, val.Prompt)

		jsonString, _ := json.Marshal(val.Answer)
		_, err = stmt.Exec(key, val.Answer.Properties.Name, val.Prompt, string(jsonString))
		if err != nil {
			log.Fatal(err)
		}
	}
}
