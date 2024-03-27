package service_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gpMocks "github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer/mocks"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/server"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_service_GetStations(t *testing.T) {
	type fields struct {
		uc *gpMocks.MockUseCase
	}
	tests := []struct {
		name     string
		setup    func(f fields)
		wantCode int
		want     map[string]any
	}{
		{
			name: "get stations happy",
			setup: func(f fields) {
				f.uc.EXPECT().
					GetStations(mock.Anything).
					Return(
						[]*models.Station{{Name: "a", Slug: "a", StreamUrl: "a"}},
						nil,
					)
			},
			wantCode: http.StatusOK,
			want:     map[string]any{"stations": []*models.Station{{Name: "a", Slug: "a", StreamUrl: "a"}}},
		},
		{
			name: "get stations unhappy",
			setup: func(f fields) {
				f.uc.EXPECT().
					GetStations(mock.Anything).
					Return(nil, errors.New("boom"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		f := fields{
			uc: gpMocks.NewMockUseCase(t),
		}
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(f)
			}

			router := givenService(f.uc)

			w := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/stations", nil)
			assert.NoError(t, err)

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.want != nil {
				wantj, _ := json.Marshal(tt.want)
				assert.Equal(t, string(wantj), w.Body.String())
			}
		})
	}
}

func givenService(uc *gpMocks.MockUseCase) *gin.Engine {
	s := service.NewService(uc)
	srv := server.NewServer(s)
	return srv.Router
}
