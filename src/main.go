package src

import (
	"embed"

	"seanmolenaar.eu/de/src/db"
	"seanmolenaar.eu/de/src/ui"
)

func WipeDB() { db.WipeDB() }
func SeedDB() { db.SeedDB() }

func SetAssetFS(fs embed.FS) { ui.AssetFS = fs }

func SetHtmlFS(fs embed.FS) { ui.HtmlFS = fs }
