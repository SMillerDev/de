package src

import (
	"embed"
	"log"

	"seanmolenaar.eu/de/src/db"
	"seanmolenaar.eu/de/src/ui"
)

func WipeDB() { db.WipeDB() }
func SeedDB() { db.SeedDB() }

func InitDB() {
	dbVar, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	db.Db = dbVar
}

func DbClose() {
	db.Db.Close()
}

func SetAssetFS(fs embed.FS) { ui.AssetFS = fs }

func SetHtmlFS(fs embed.FS) { ui.HtmlFS = fs }
