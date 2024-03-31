package models

import "github.com/jj-style/gobal-player/pkg/globalplayer/models/nextjs"

type Station struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	StreamUrl string `json:"streamUrl"`
	Id        string `json:"id"`
	Tagline   string `json:"tagline"`
	ImageUrl  string `json:"imageUrl"`
}

func StationFromApiModel(s *nextjs.StationBrand) *Station {
	return &Station{
		Name:      s.Name,
		Slug:      s.Slug,
		StreamUrl: s.NationalStation.StreamURL,
		Id:        s.ID,
		Tagline:   s.NationalStation.Tagline,
		ImageUrl:  s.LogoUnstackedWithoutBackground,
	}
}
