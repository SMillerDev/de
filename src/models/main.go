package models

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type GeoJSONAnswer struct {
	Type       string `json:"type,case:ignore"`
	Properties struct {
		Name string `json:"name,case:ignore"`
	} `json:"properties,case:ignore"`
	Geometry struct {
		Type        string        `json:"type,case:ignore"`
		Coordinates [][][]float32 `json:"coordinates,case:ignore"`
	} `json:"geometry,case:ignore"`
}

type GameItem struct {
	Id      string            `json:"-"`
	Prompts map[string]string `json:"prompts,case:ignore"`
	PlayAt  time.Time         `json:"-"`
	Answer  GeoJSONAnswer     `json:"answer,case:ignore"`
	Wiki   map[string]string `json:"wikis,case:ignore"`
	Tags []string`json:"tags,case:ignore"`
}

type TodayApiResponse struct {
	PlayDate string              `json:"playDate,case:ignore"`
	Data     map[string]GameItem `json:"data,case:ignore"`
}

type Coordinate struct {
	Lattitude float64
	Longitude float64
}

type ReportValues struct {
	Id         string
	Coordinate Coordinate
}

func (item GameItem) NewItemFromDbEntry(id string, name string, playAtString string, promptBytes []byte, answerBytes []byte) GameItem {
	var answer GeoJSONAnswer
	json.Unmarshal(answerBytes, &answer)
	answer.Properties.Name = name

	var prompts map[string]string
	err := json.Unmarshal(promptBytes, &prompts)
	if err != nil {
		log.Fatal(err)
	}

	if playAtString != "" {
		playAt, _ := time.Parse(time.DateOnly, playAtString)
		item.PlayAt = playAt
	} else {
		item.PlayAt = time.Time{}
	}

	item.Id = id
	item.Prompts = prompts
	item.Answer = answer

	return item
}

func (item ReportValues) NewItemFromDbEntry(id string, latString string, longString string) ReportValues {
	item.Id = id
	lat, _ := strconv.ParseFloat(latString, 32)
	long, _ := strconv.ParseFloat(longString, 32)
	item.Coordinate = Coordinate{
		Lattitude: lat,
		Longitude: long,
	}

	return item
}
