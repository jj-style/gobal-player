package models

import (
	"time"
)

type Episode struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ImageUrl        string    `json:"imageUrl"`
	StreamUrl       string    `json:"streamUrl"`
	Aired           time.Time `json:"aired"`
	Until           time.Time `json:"until"`
	Duration        string    `json:"duration"`
	DurationSeconds int       `json:"durationSeconds"`
	Availability    string    `json:"availability"`
}
