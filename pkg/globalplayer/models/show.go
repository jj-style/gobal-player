package models

import "github.com/jj-style/gobal-player/pkg/globalplayer/models/nextjs"

type Show struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

func ShowFromApiModel(s *nextjs.CatchupInfo) *Show {
	return &Show{
		Id:       s.ID,
		Name:     s.Title,
		ImageUrl: s.ImageURL,
	}
}
