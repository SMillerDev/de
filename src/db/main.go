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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS promptsNew (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
	    playAt TEXT NOT NULL DEFAULT '1970-01-01',
		prompt TEXT NOT NULL DEFAULT '{}' CHECK(JSON_VALID(prompt)),
		wiki TEXT NOT NULL DEFAULT '{}' CHECK(JSON_VALID(wiki)),
		tags TEXT NOT NULL DEFAULT '[]' CHECK(JSON_VALID(tags)),
	 	answer TEXT UNIQUE NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new prompts: %s", err.Error())
	}

	// Migration
	_, err = db.Exec(`REPLACE INTO promptsNew SELECT id,name,playAt,JSON_OBJECT('en', prompt),NULL,NULL, answer FROM prompts  WHERE json_valid(prompt) = 0;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to migrate old prompts: %s", err.Error())
	}

	// Keep items
	_, err = db.Exec(`REPLACE INTO promptsNew SELECT id,name,playAt,prompt,NULL,NULL,answer FROM prompts WHERE json_valid(prompt) = 1;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to recreate prompts: %s", err.Error())
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS prompts;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to drop old prompts: %s", err.Error())
	}

	_, err = db.Exec(`ALTER TABLE promptsNew RENAME TO prompts;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to rename new prompts: %s", err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reports (id NOT NULL, lat REAL NOT NULL, long REAL NOT NULL, FOREIGN KEY(id) REFERENCES prompts(id))`)
	if err != nil {
		return nil, fmt.Errorf("Failed to create reports: %s", err.Error())
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
