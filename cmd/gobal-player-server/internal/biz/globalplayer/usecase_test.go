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
		want    []*models.Station
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetStations().
					Return([]*models.Station{
						{
							Id:        "id",
							Name:      "name",
							Slug:      "slug",
							StreamUrl: "url",
						},
					}, nil)
			},
			want: []*models.Station{
				{Id: "id", Name: "name", Slug: "slug", StreamUrl: "url"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "unhappy",
			setup: func(f fields, _ args) {
				f.gp.EXPECT().GetStations().Return(nil, errors.New("boom"))
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
