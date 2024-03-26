package globalplayer

import "github.com/jj-style/gobal-player/pkg/globalplayer/models"

type Station struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	StreamUrl string `json:"streamUrl"`
	Id        string `json:"id"`
}

func StationFromApiModel(s models.StationBrand) *Station {
	return &Station{
		Name:      s.Name,
		Slug:      s.Slug,
		StreamUrl: s.NationalStation.StreamURL,
		Id:        s.ID,
	}
}
