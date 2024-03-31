package service_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
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
			request, err := http.NewRequest(http.MethodGet, "/api/stations", nil)
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

func Test_service_GetShows(t *testing.T) {
	type fields struct {
		uc *gpMocks.MockUseCase
	}
	type args struct {
		slug string
	}
	tests := []struct {
		name     string
		args     args
		setup    func(f fields)
		wantCode int
		want     map[string]any
	}{
		{
			name: "happy",
			args: args{slug: "slug"},
			setup: func(f fields) {
				f.uc.EXPECT().
					GetShows(mock.Anything, "slug").
					Return(
						[]*models.Show{{Name: "name"}},
						nil,
					)
			},
			wantCode: http.StatusOK,
			want:     map[string]any{"shows": []*models.Show{{Name: "name"}}},
		},
		{
			name: "unhappy",
			args: args{slug: "slug"},
			setup: func(f fields) {
				f.uc.EXPECT().
					GetShows(mock.Anything, "slug").
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
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/shows/%s", tt.args.slug), nil)
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

func Test_service_GetEpisodes(t *testing.T) {
	type fields struct {
		uc *gpMocks.MockUseCase
	}
	type args struct {
		slug string
		id   string
	}
	tests := []struct {
		name     string
		args     args
		setup    func(f fields)
		wantCode int
		want     map[string]any
	}{
		{
			name: "happy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields) {
				f.uc.EXPECT().
					GetEpisodes(mock.Anything, "slug", "id").
					Return(
						[]*models.Episode{{Name: "name"}},
						nil,
					)
			},
			wantCode: http.StatusOK,
			want:     map[string]any{"episodes": []*models.Episode{{Name: "name"}}},
		},
		{
			name: "unhappy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields) {
				f.uc.EXPECT().
					GetEpisodes(mock.Anything, "slug", "id").
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
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/episodes/%s/%s", tt.args.slug, tt.args.id), nil)
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

func Test_service_GetEpisodesRss(t *testing.T) {
	type fields struct {
		uc *gpMocks.MockUseCase
	}
	type args struct {
		slug string
		id   string
	}
	tests := []struct {
		name            string
		args            args
		setup           func(f fields)
		wantCode        int
		wantContentType string
	}{
		{
			name: "happy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields) {
				f.uc.EXPECT().GetEpisodesFeed(mock.Anything, "slug", "id").
					Return(&feeds.Feed{Title: "title", Description: "description", Id: "id"}, nil)
			},
			wantCode:        http.StatusOK,
			wantContentType: "application/xml",
		},
		{
			name: "unhappy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields) {
				f.uc.EXPECT().GetEpisodesFeed(mock.Anything, "slug", "id").
					Return(nil, errors.New("boom"))
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "application/json; charset=utf-8",
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
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/episodes/%s/%s/rss", tt.args.slug, tt.args.id), nil)
			assert.NoError(t, err)

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantContentType, w.Header().Get("content-type"))
		})
	}
}

func Test_service_GetAllShowsRss(t *testing.T) {
	type fields struct {
		uc *gpMocks.MockUseCase
	}
	type args struct {
		slug string
	}
	tests := []struct {
		name            string
		args            args
		setup           func(f fields)
		wantCode        int
		wantContentType string
	}{
		{
			name: "happy",
			args: args{slug: "slug"},
			setup: func(f fields) {
				f.uc.EXPECT().GetAllShowsFeed(mock.Anything, "slug").
					Return(&feeds.Feed{Title: "title", Description: "description", Id: "id"}, nil)
			},
			wantCode:        http.StatusOK,
			wantContentType: "application/xml",
		},
		{
			name: "unhappy",
			args: args{slug: "slug"},
			setup: func(f fields) {
				f.uc.EXPECT().GetAllShowsFeed(mock.Anything, "slug").
					Return(nil, errors.New("boom"))
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "application/json; charset=utf-8",
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
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/shows/%s/rss", tt.args.slug), nil)
			assert.NoError(t, err)

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantContentType, w.Header().Get("content-type"))
		})
	}
}

func givenService(uc *gpMocks.MockUseCase) *gin.Engine {
	s := service.NewService(uc)
	srv := server.NewServer(s)
	return srv.Router
}
