package globalplayer_test

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/gorilla/feeds"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
	gpMocks "github.com/jj-style/gobal-player/pkg/globalplayer/mocks"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	restyMocks "github.com/jj-style/gobal-player/pkg/resty/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_useCase_GetStations(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
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
						{Id: "id", Name: "name"},
					}, nil)
			},
			want: []*models.Station{
				{Id: "id", Name: "name"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "unhappy",
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetStations().
					Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				gp: gpMocks.NewMockGlobalPlayer(t),
				hc: restyMocks.NewMockHttpClient(t),
			}

			if tt.setup != nil {
				tt.setup(f, tt.args)
			}

			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetStations(ctx)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetShows(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		slug string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields, args)
		want    []*models.Show
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			args: args{slug: "slug"},
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetShows("slug").
					Return([]*models.Show{
						{Id: "id", Name: "name"},
					}, nil)
			},
			want: []*models.Show{
				{Id: "id", Name: "name"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "unhappy",
			args: args{slug: "slug"},
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetShows("slug").
					Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				gp: gpMocks.NewMockGlobalPlayer(t),
				hc: restyMocks.NewMockHttpClient(t),
			}

			if tt.setup != nil {
				tt.setup(f, tt.args)
			}

			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetShows(ctx, tt.args.slug)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetEpisodes(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		slug string
		id   string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields, args)
		want    []*models.Episode
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetEpisodes("slug", "id").
					Return([]*models.Episode{
						{Id: "id", Name: "name"},
					}, nil)
			},
			want: []*models.Episode{
				{Id: "id", Name: "name"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "unhappy",
			args: args{slug: "slug", id: "id"},
			setup: func(f fields, _ args) {
				f.gp.EXPECT().
					GetEpisodes("slug", "id").
					Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				gp: gpMocks.NewMockGlobalPlayer(t),
				hc: restyMocks.NewMockHttpClient(t),
			}

			if tt.setup != nil {
				tt.setup(f, tt.args)
			}

			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetEpisodes(ctx, tt.args.slug, tt.args.id)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetStation(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		stationSlug string
	}
	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		want    *models.Station
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetStations().
					Return([]*models.Station{{Id: "a", Name: "a", Slug: "a"}, {Id: "b", Name: "b", Slug: "b"}}, nil)
			},
			args:    args{stationSlug: "b"},
			want:    &models.Station{Id: "b", Name: "b", Slug: "b"},
			wantErr: assert.NoError,
		},
		{
			name: "error getting stations",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetStations().
					Return(nil, errors.New("boom"))
			},
			args:    args{stationSlug: "a"},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "station not found",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetStations().
					Return([]*models.Station{{Id: "a", Name: "a", Slug: "a"}}, nil)
			},
			args:    args{stationSlug: "b"},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		f := &fields{
			gp: gpMocks.NewMockGlobalPlayer(t),
			hc: restyMocks.NewMockHttpClient(t),
		}
		if tt.setup != nil {
			tt.setup(f)
		}
		t.Run(tt.name, func(t *testing.T) {
			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetStation(ctx, tt.args.stationSlug)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetShow(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		stationSlug string
		showId      string
	}
	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		want    *models.Show
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetShows("station").
					Return([]*models.Show{{Id: "a", Name: "a"}, {Id: "b", Name: "b"}}, nil)
			},
			args:    args{stationSlug: "station", showId: "b"},
			want:    &models.Show{Id: "b", Name: "b"},
			wantErr: assert.NoError,
		},
		{
			name: "error getting shows",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetShows("station").
					Return(nil, errors.New("boom"))
			},
			args:    args{stationSlug: "station", showId: "b"},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "show not found",
			setup: func(f *fields) {
				f.gp.EXPECT().
					GetShows("station").
					Return([]*models.Show{{Id: "a", Name: "a"}}, nil)
			},
			args:    args{stationSlug: "station", showId: "b"},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		f := &fields{
			gp: gpMocks.NewMockGlobalPlayer(t),
			hc: restyMocks.NewMockHttpClient(t),
		}
		if tt.setup != nil {
			tt.setup(f)
		}
		t.Run(tt.name, func(t *testing.T) {
			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetShow(ctx, tt.args.stationSlug, tt.args.showId)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetEpisodesFeed(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		stationSlug string
		showId      string
	}
	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		want    *feeds.Feed
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			args: args{stationSlug: "station", showId: "show"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetEpisodes("station", "show").
					Return([]*models.Episode{{Id: "id", Name: "episode 1", Description: "episode", StreamUrl: "episode.mp3", Duration: "00:30:00"}}, nil)

				f.gp.EXPECT().
					GetShows("station").
					Return([]*models.Show{{Id: "show", Name: "show", ImageUrl: "show.jpg"}}, nil)

				expectReq, _ := http.NewRequest(http.MethodHead, "episode.mp3", nil)
				f.hc.EXPECT().Do(expectReq).Return(&http.Response{ContentLength: 100}, nil)
			},
			want: &feeds.Feed{
				Title:       "show",
				Description: "episode",
				Subtitle:    "episode",
				Image:       &feeds.Image{Url: "show.jpg"},
				ITunes: &feeds.ITunesFeed{
					Explicit: false,
					Type:     feeds.ITunesFeedTypeEpisodic,
					Title:    "show",
					Image:    &feeds.ITunesImage{Href: "show.jpg"},
				},
				Items: []*feeds.Item{
					{
						Id:          "id",
						Title:       "episode 1: Monday 01 January 0001",
						Description: "episode<br/><br/>Available until Monday 01 January 0001 00:00:00.",
						Enclosure:   &feeds.Enclosure{Url: "episode.mp3", Type: "audio/mpeg", Length: "100"},
						Link:        &feeds.Link{Href: "episode.mp3"},
						ITunes: &feeds.ITunesItem{
							Duration:    "30m0s",
							EpisodeType: feeds.ITunesEpisodeTypeFull,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error getting episodes",
			args: args{stationSlug: "station", showId: "show"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetEpisodes("station", "show").
					Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error getting show",
			args: args{stationSlug: "station", showId: "show"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetEpisodes("station", "show").
					Return([]*models.Episode{{Id: "id", Name: "episode 1", Description: "episode", StreamUrl: "episode.mp3"}}, nil)

				f.gp.EXPECT().
					GetShows("station").
					Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		f := &fields{
			gp: gpMocks.NewMockGlobalPlayer(t),
			hc: restyMocks.NewMockHttpClient(t),
		}
		if tt.setup != nil {
			tt.setup(f)
		}
		t.Run(tt.name, func(t *testing.T) {
			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetEpisodesFeed(ctx, tt.args.stationSlug, tt.args.showId)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_useCase_GetAllShowsFeed(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		gp *gpMocks.MockGlobalPlayer
		hc *restyMocks.MockHttpClient
	}
	type args struct {
		stationSlug string
	}
	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		want    *feeds.Feed
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			args: args{stationSlug: "station"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetStations().Return([]*models.Station{{Name: "station", Slug: "station", Tagline: "a cool station", ImageUrl: "station.jpg"}}, nil)

				f.gp.EXPECT().
					GetShows("station").
					Return([]*models.Show{
						{Id: "show1", Name: "show1"},
						{Id: "show2", Name: "show2"},
					}, nil)

				f.gp.EXPECT().GetEpisodes("station", "show1").
					Return([]*models.Episode{{Id: "show1id1", Name: "show 1 episode 1", Description: "show 1 episode 1", StreamUrl: "s1ep1.mp3", Duration: "00:30:00"}}, nil)

				f.gp.EXPECT().GetEpisodes("station", "show2").
					Return([]*models.Episode{{Id: "show2id1", Name: "show 2 episode 1", Description: "show 2 episode 1", StreamUrl: "s2ep1.mp3", Duration: "00:30:00"}}, nil)

				expectReq1, _ := http.NewRequest(http.MethodHead, "s1ep1.mp3", nil)
				f.hc.EXPECT().Do(expectReq1).Return(&http.Response{ContentLength: 100}, nil)
				expectReq2, _ := http.NewRequest(http.MethodHead, "s2ep1.mp3", nil)
				f.hc.EXPECT().Do(expectReq2).Return(&http.Response{ContentLength: 200}, nil)

			},
			want: &feeds.Feed{
				Title:       "station",
				Description: "a cool station",
				Subtitle:    "a cool station",
				Image:       &feeds.Image{Url: "station.jpg"},
				ITunes: &feeds.ITunesFeed{
					Image:    &feeds.ITunesImage{Href: "station.jpg"},
					Explicit: false,
					Type:     feeds.ITunesFeedTypeEpisodic,
					Title:    "station",
				},
				Items: []*feeds.Item{
					{
						Id:          "show1id1",
						Title:       "show 1 episode 1: Monday 01 January 0001",
						Description: "show 1 episode 1<br/><br/>Available until Monday 01 January 0001 00:00:00.",
						Enclosure:   &feeds.Enclosure{Url: "s1ep1.mp3", Type: "audio/mpeg", Length: "100"},
						Link:        &feeds.Link{Href: "s1ep1.mp3"},
						ITunes: &feeds.ITunesItem{
							Duration:    "30m0s",
							EpisodeType: feeds.ITunesEpisodeTypeFull,
						},
					},
					{
						Id:          "show2id1",
						Title:       "show 2 episode 1: Monday 01 January 0001",
						Description: "show 2 episode 1<br/><br/>Available until Monday 01 January 0001 00:00:00.",
						Enclosure:   &feeds.Enclosure{Url: "s2ep1.mp3", Type: "audio/mpeg", Length: "200"},
						Link:        &feeds.Link{Href: "s2ep1.mp3"},
						ITunes: &feeds.ITunesItem{
							Duration:    "30m0s",
							EpisodeType: feeds.ITunesEpisodeTypeFull,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error getting episodes",
			args: args{stationSlug: "station"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetStations().Return([]*models.Station{{Name: "station", Slug: "station", Tagline: "a cool station", ImageUrl: "station.jpg"}}, nil)

				f.gp.EXPECT().
					GetShows("station").
					Return([]*models.Show{
						{Id: "show1", Name: "show1"},
						{Id: "show2", Name: "show2"},
					}, nil)

				f.gp.EXPECT().GetEpisodes("station", "show1").
					Return(nil, errors.New("boom"))

				f.gp.EXPECT().GetEpisodes("station", "show2").
					Return([]*models.Episode{{Id: "show2id1", Name: "show 2 episode 1", Description: "show 2 episode 1", StreamUrl: "s2ep1.mp3"}}, nil).Maybe()

			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error getting shows",
			args: args{stationSlug: "station"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetStations().Return([]*models.Station{{Name: "station", Slug: "station", Tagline: "a cool station", ImageUrl: "station.jpg"}}, nil)

				f.gp.EXPECT().
					GetShows("station").
					Return(nil, errors.New("boom"))

			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error getting station",
			args: args{stationSlug: "station"},
			setup: func(f *fields) {
				f.gp.EXPECT().GetStations().Return(nil, errors.New("boom"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		f := &fields{
			gp: gpMocks.NewMockGlobalPlayer(t),
			hc: restyMocks.NewMockHttpClient(t),
		}
		if tt.setup != nil {
			tt.setup(f)
		}
		t.Run(tt.name, func(t *testing.T) {
			u := globalplayer.NewUseCase(f.gp, f.hc)
			got, err := u.GetAllShowsFeed(ctx, tt.args.stationSlug)
			tt.wantErr(t, err)
			if tt.want != nil {

				slices.SortFunc(got.Items, func(a, b *feeds.Item) int { return strings.Compare(a.Id, b.Id) })
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
