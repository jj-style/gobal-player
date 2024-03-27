package globalplayer

import (
	"context"

	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
)

type UseCase interface {
	GetStations(ctx context.Context) ([]*models.Station, error)
}

type useCase struct {
	gp globalplayer.GlobalPlayer
}

func (u *useCase) GetStations(ctx context.Context) ([]*models.Station, error) {
	return u.gp.GetStations()
}

func NewUseCase(gp globalplayer.GlobalPlayer) UseCase {
	return &useCase{
		gp: gp,
	}
}
