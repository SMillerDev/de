package db

import (
	"database/sql"
	_ "embed"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func GetItemById(id string) ([]models.GameItem, error) {
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		return []models.GameItem{}, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, playAt, prompt, answer FROM prompts WHERE id = ?;", id)
	if err != nil {
		return []models.GameItem{}, err
	}
	var result = []models.GameItem{}
	defer rows.Close()
	for rows.Next() {
		var name string
		var playAtNil sql.NullTime
		var prompt string
		var answerJson []byte
		err = rows.Scan(&name, &playAtNil, &prompt, &answerJson)
		if err != nil {
			log.Fatal(err)
		}

		playAt := time.Time{}
		if playAtNil.Valid {
			playAt = playAtNil.Time
		}

		result = append(result, NewItemFromDbEntry(id, name, playAt, prompt, answerJson))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func GetDBItems(limit int) []models.GameItem {
	db, err := sql.Open("sqlite3", getDBPath())
	var result = []models.GameItem{}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, name, playAt, prompt, answer FROM prompts ORDER BY id LIMIT ?", limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var playAtNil sql.NullTime
		var prompt string
		var answerJson []byte
		err = rows.Scan(&id, &name, &playAtNil, &prompt, &answerJson)
		if err != nil {
			log.Fatal(err)
		}
		playAt := time.Time{}
		if playAtNil.Valid {
			playAt = playAtNil.Time
		}

		result = append(result, NewItemFromDbEntry(id, name, playAt, prompt, answerJson))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
