package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS prompts (id PRIMARY KEY, playAt, name NOT NULL, prompt NOT NULL, answer json UNIQUE NOT NULL)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reports (id NOT NULL, lat REAL NOT NULL, long REAL NOT NULL, FOREIGN KEY(id) REFERENCES prompts(id))`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

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
