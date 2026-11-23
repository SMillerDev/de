package db

import (
	"database/sql"
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func GetItemById(id string) ([]models.GameItem, error) {
	rows, err := Db.Query("SELECT name, playAt, prompt, answer FROM prompts WHERE id = ?;", id)
	if err != nil {
		return nil, err
	}
	var result = []models.GameItem{}
	defer rows.Close()
	for rows.Next() {
		var name string
		var playAtNull sql.NullString
		var promptsJson []byte
		var answerJson []byte
		err = rows.Scan(&name, &playAtNull, &promptsJson, &answerJson)
		if err != nil {
			return nil, err
		}

		playAtString := ""
		if playAtNull.Valid {
			playAtString = playAtNull.String
		}

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, playAtString, promptsJson, answerJson))
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetDBItems(limit int) []models.GameItem {
	var result = []models.GameItem{}
	rows, err := Db.Query("SELECT id, name, playAt, prompt, answer FROM prompts ORDER BY playAt DESC LIMIT ?", limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var playAtNull sql.NullString
		var promptsJson []byte
		var answerJson []byte
		err = rows.Scan(&id, &name, &playAtNull, &promptsJson, &answerJson)
		if err != nil {
			log.Fatal(err)
		}

		playAtString := ""
		if playAtNull.Valid {
			playAtString = playAtNull.String
		}

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, playAtString, promptsJson, answerJson))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func DeleteItemById(id string) error {
	_, err := Db.Exec("DELETE FROM prompts WHERE id = ?;", id)
	if err != nil {
		return err
	}
	return nil
}

func SaveItem(item models.GameItem) error {
	stmt, err := Db.Prepare("REPLACE INTO prompts(id, playAt, name, prompt, answer) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Printf("ID: %s, Prompt: %s\n", item.Id, item.Prompts["en"])

	jsonAnswerString, err := json.Marshal(item.Answer)
	if err != nil {
		return err
	}
	jsonPromptString, err := json.Marshal(item.Prompts)
	if err != nil {
		return err
	}
	var sqlTime any
	if item.PlayAt.IsZero() {
		sqlTime = nil
	} else {
		sqlTime = item.PlayAt.Format(time.DateOnly)
	}
	if len(item.Prompts) == 0 {
		return fmt.Errorf("Empty value for prompt!")
	}
	if item.Answer.Properties.Name == "" {
		return fmt.Errorf("Empty value for name!")
	}

	_, err = stmt.Exec(item.Id, sqlTime, item.Answer.Properties.Name, string(jsonPromptString), string(jsonAnswerString))
	return err
}

func GetReports(id string, limit int) []models.ReportValues {
	var result = []models.ReportValues{}
	rows, err := Db.Query("SELECT id, lat, long FROM reports WHERE id=? LIMIT ?", id, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var lat string
		var long string
		err = rows.Scan(&id, &lat, &long)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, models.ReportValues{}.NewItemFromDbEntry(id, lat, long))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
