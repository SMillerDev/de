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

func Report(report models.ReportValues) error {
	stmt, err := Db.Prepare("INSERT INTO reports(id, lat, long) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(report.Id, report.Coordinate.Lattitude, report.Coordinate.Longitude)
	return err
}
