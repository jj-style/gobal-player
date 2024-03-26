package globalplayer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
	gpMocks "github.com/jj-style/gobal-player/pkg/globalplayer/mocks"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/stretchr/testify/assert"
)

func Test_useCase_GetStations(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
	}
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields, args)
		want    []globalplayer.Station
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetStations().
					Return(newStationsPageResponse(
						models.StationBrand{
							Slug: "slug",
							Name: "name",
							NationalStation: models.NationalStation{
								StreamURL: "url",
							},
							ID: "id",
						}),
						nil)
			},
			want: []globalplayer.Station{
				{Name: "name", Slug: "slug", StreamUrl: "url", Id: "id"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "unhappy",
			setup: func(f fields, _ args) {
				f.gp.EXPECT().GetStations().Return(models.StationsPageResponse{}, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				gp: gpMocks.NewMockGlobalPlayer(t),
			}

			if tt.setup != nil {
				tt.setup(f, tt.args)
			}

			u := globalplayer.NewUseCase(f.gp)
			got, err := u.GetStations(ctx)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// test helper to return the StationBrands
func newStationsPageResponse(stations ...models.StationBrand) models.StationsPageResponse {
	return models.StationsPageResponse{
		PageProps: models.StatsionsPageProps{
			Feature: models.Feature{
				Blocks: []models.Block{
					models.Block{
						Brands: stations,
					},
				},
			},
		},
	}
}
