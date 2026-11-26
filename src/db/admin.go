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
	rows, err := Db.Query("SELECT name, playAt, prompt, wiki, tags, answer FROM prompts WHERE id = ?;", id)
	if err != nil {
		return nil, err
	}
	var result = []models.GameItem{}
	defer rows.Close()
	for rows.Next() {
		var name string
		var playAtNull sql.NullString
		var promptsJson []byte
		var wikiJson []byte
		var tagsJson []byte
		var answerJson []byte
		err = rows.Scan(&name, &playAtNull, &promptsJson, &wikiJson, &tagsJson, &answerJson)
		if err != nil {
			return nil, err
		}

		playAtString := ""
		if playAtNull.Valid {
			playAtString = playAtNull.String
		}

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, playAtString, promptsJson, tagsJson, wikiJson, answerJson))
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetDBItems(limit int) []models.GameItem {
	var result = []models.GameItem{}
	rows, err := Db.Query("SELECT id, name, playAt, prompt, tags, wiki, answer FROM prompts ORDER BY playAt DESC LIMIT ?", limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var playAtNull sql.NullString
		var promptsJson []byte
		var tagJson []byte
		var wikiJson []byte
		var answerJson []byte
		err = rows.Scan(&id, &name, &playAtNull, &promptsJson, &tagJson, &wikiJson, &answerJson)
		if err != nil {
			log.Fatal(err)
		}

		playAtString := ""
		if playAtNull.Valid {
			playAtString = playAtNull.String
		}

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, playAtString, promptsJson, tagJson, wikiJson, answerJson))
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

func toJson(in any) (string, error) {
	jsonString, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

func SaveItem(item models.GameItem) error {
	stmt, err := Db.Prepare("REPLACE INTO prompts(id, playAt, name, prompt, wiki, tags, answer) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Printf("ID: %s, Prompt: %s\n", item.Id, item.Prompts["en"])

	answerString, err := toJson(item.Answer)
	if err != nil {
		return err
	}
	promptString, err := toJson(item.Prompts)
	if err != nil {
		return err
	}
	wikiString, err := toJson(item.Wiki)
	if err != nil {
		return err
	}
	tagString, err := toJson(item.Tags)
	if err != nil {
		return err
	}
	var sqlTime any
	if item.PlayAt.IsZero() {
		sqlTime = nil
	} else {
		sqlTime = item.PlayAt.Format(time.DateOnly)
	}
	if item.Answer.Properties.Name == "" {
		return fmt.Errorf("Empty value for name!")
	}

	_, err = stmt.Exec(item.Id, sqlTime, item.Answer.Properties.Name, promptString, wikiString, tagString, answerString)
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

func GetDateCounts() map[string]int {
	var result = make(map[string]int)
	rows, err := Db.Query("SELECT playAt, count(id) AS count FROM prompts GROUP BY playAt;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		var count int
		err = rows.Scan(&date, &count)
		if err != nil {
			log.Fatal(err)
		}

		result[date] = count
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetTagCounts() map[string]int {
	var result = make(map[string]int)
	rows, err := Db.Query("SELECT tag.value AS tag, COUNT(tag.value) AS count FROM prompts JOIN json_each(prompts.tags) AS tag GROUP BY tag.value;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		var count int
		err = rows.Scan(&tag, &count)
		if err != nil {
			log.Fatal(err)
		}

		result[tag] = count
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
