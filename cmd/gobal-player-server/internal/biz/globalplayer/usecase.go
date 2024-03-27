package globalplayer

import (
	"context"

	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
)

type UseCase interface {
	GetStations(context.Context) ([]*models.Station, error)
	GetShows(context.Context, string) ([]*models.Show, error)
	GetEpisodes(context.Context, string, string) ([]*models.Episode, error)
}

type useCase struct {
	gp globalplayer.GlobalPlayer
}

func (u *useCase) GetStations(ctx context.Context) ([]*models.Station, error) {
	return u.gp.GetStations()
}

func (u *useCase) GetShows(ctx context.Context, stationSlug string) ([]*models.Show, error) {
	return u.gp.GetShows(stationSlug)
}

func (u *useCase) GetEpisodes(ctx context.Context, stationSlug, showId string) ([]*models.Episode, error) {
	return u.gp.GetEpisodes(stationSlug, showId)
}

func NewUseCase(gp globalplayer.GlobalPlayer) UseCase {
	return &useCase{
		gp: gp,
	}
}
