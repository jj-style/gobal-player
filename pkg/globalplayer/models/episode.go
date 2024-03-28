package models

import (
	"time"

	"github.com/jj-style/gobal-player/pkg/globalplayer/models/nextjs"
)

type Episode struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ImageUrl     string    `json:"imageUrl"`
	StreamUrl    string    `json:"streamUrl"`
	Aired        time.Time `json:"aired"`
	Until        time.Time `json:"until"`
	Duration     string    `json:"duration"`
	Availability string    `json:"availability"`
}

func EpisodeFromApiModel(s *nextjs.Episode) *Episode {
	return &Episode{
		Id:           s.ID,
		Name:         s.Title,
		Description:  s.Description,
		ImageUrl:     s.ImageURL,
		StreamUrl:    s.StreamURL,
		Aired:        s.StartDate,
		Until:        s.AvailableUntil,
		Duration:     s.Duration,
		Availability: s.Availability,
	}
}
