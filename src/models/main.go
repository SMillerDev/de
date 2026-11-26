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
	Wiki    map[string]string `json:"wikis,case:ignore"`
	Tags    []string          `json:"tags,case:ignore"`
	Reports []ReportValues    `json:"-"`
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

func unmarshalMap(in []byte) (map[string]string, error) {
	var out map[string]string
	err := json.Unmarshal(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (item GameItem) NewItemFromDbEntry(id string, name string, playAtString string, promptBytes []byte, tagBytes []byte, wikiBytes []byte, answerBytes []byte) GameItem {
	var answer GeoJSONAnswer
	json.Unmarshal(answerBytes, &answer)
	answer.Properties.Name = name

	prompts, err := unmarshalMap(promptBytes)
	if err != nil {
		log.Fatal(err)
	}

	wiki, err := unmarshalMap(wikiBytes)
	if err != nil {
		log.Fatal(err)
	}

	var tags []string
	err = json.Unmarshal(tagBytes, &tags)
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
	item.Tags = tags
	item.Wiki = wiki

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
