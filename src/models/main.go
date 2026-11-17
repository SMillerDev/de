package models

import "time"

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
	Id     string        `json:"-"`
	Prompt string        `json:"prompt,case:ignore"`
	PlayAt time.Time     `json:"-"`
	Answer GeoJSONAnswer `json:"answer,case:ignore"`
}

type TodayApiResponse struct {
	Data map[string]GameItem `json:"data,case:ignore"`
}

type ReportValues struct {
	Id        string
	Lattitude float64
	Longitude float64
}
