package db

import (
	_ "embed"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"seanmolenaar.eu/de/src/models"
)

func GetItemsToday(currentTime time.Time) []models.GameItem {
	var result = []models.GameItem{}
	rows, err := Db.Query("SELECT id, name, prompt, tags, wiki, answer FROM prompts WHERE playAt=? ORDER BY id ", currentTime.Format(time.DateOnly))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var promptsJson []byte
		var tagJson []byte
		var wikiJson []byte
		var answerJson []byte
		err = rows.Scan(&id, &name, &promptsJson, &tagJson, &wikiJson, &answerJson)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, models.GameItem{}.NewItemFromDbEntry(id, name, currentTime.Format(time.DateOnly), promptsJson, tagJson, wikiJson, answerJson, 0))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func Report(report models.ReportValues) error {
	stmt, err := Db.Prepare("INSERT INTO reports(id, lat, long) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(report.Id, report.Coordinate.Lattitude, report.Coordinate.Longitude)
	return err
}
