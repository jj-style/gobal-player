package models

import "github.com/jj-style/gobal-player/pkg/globalplayer/models/nextjs"

type LiveStation struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	StreamUrl string `json:"streamUrl"`
	Id        string `json:"id"`
	Tagline   string `json:"tagline"`
	ImageUrl  string `json:"imageUrl"`
}

func LiveStationFromApiModel(b *nextjs.Brand) *LiveStation {
	return &LiveStation{
		Id:        b.BrandID,
		Name:      b.BrandName,
		ImageUrl:  b.BrandLogo,
		StreamUrl: b.StreamURL,
		Tagline:   b.Tagline,
		Slug:      b.BrandSlug,
	}
}
