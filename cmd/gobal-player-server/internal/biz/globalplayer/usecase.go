package globalplayer

import (
	"context"

	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/samber/lo"
)

type UseCase interface {
	GetStations(ctx context.Context) ([]Station, error)
}

type useCase struct {
	gp globalplayer.GlobalPlayer
}

func (u *useCase) GetStations(ctx context.Context) ([]Station, error) {
	got, err := u.gp.GetStations()
	if err != nil {
		return nil, err
	}

	brands := got.PageProps.Feature.Blocks[0].Brands

	return lo.Map(brands, func(item models.StationBrand, _ int) Station { return *StationFromApiModel(item) }), nil
}

func NewUseCase(gp globalplayer.GlobalPlayer) UseCase {
	return &useCase{
		gp: gp,
	}
}
