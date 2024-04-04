package globalplayer

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/google/wire"
	"github.com/gorilla/feeds"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	feeds2 "github.com/jj-style/gobal-player/pkg/globalplayer/feeds"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type UseCase interface {
	GetStations(context.Context) ([]*models.Station, error)
	GetStation(context.Context, string) (*models.Station, error)
	GetShows(context.Context, string) ([]*models.Show, error)
	GetShow(context.Context, string, string) (*models.Show, error)
	GetEpisodes(context.Context, string, string) ([]*models.Episode, error)
	GetEpisodesFeed(context.Context, string, string) (*feeds.Feed, error)
	GetAllShowsFeed(context.Context, string) (*feeds.Feed, error)
}

type useCase struct {
	gp globalplayer.GlobalPlayer
	hc resty.HttpClient
}

func (u *useCase) GetStations(ctx context.Context) ([]*models.Station, error) {
	return u.gp.GetStations()
}

func (u *useCase) GetStation(ctx context.Context, stationSlug string) (*models.Station, error) {
	stations, err := u.gp.GetStations()
	if err != nil {
		return nil, err
	}

	got, found := lo.Find(stations, func(item *models.Station) bool { return item.Slug == stationSlug })
	if !found {
		return nil, fmt.Errorf("station %s not found", stationSlug)
	}
	return got, nil
}

func (u *useCase) GetShows(ctx context.Context, stationSlug string) ([]*models.Show, error) {
	return u.gp.GetShows(stationSlug)
}

func (u *useCase) GetShow(ctx context.Context, stationSlug, showId string) (*models.Show, error) {
	shows, err := u.gp.GetShows(stationSlug)
	if err != nil {
		return nil, err
	}

	got, found := lo.Find(shows, func(item *models.Show) bool { return item.Id == showId })
	if !found {
		return nil, fmt.Errorf("show id=%s in station %s not found", showId, stationSlug)
	}
	return got, nil
}

func (u *useCase) GetEpisodes(ctx context.Context, stationSlug, showId string) ([]*models.Episode, error) {
	return u.gp.GetEpisodes(stationSlug, showId)
}

func (u *useCase) GetEpisodesFeed(ctx context.Context, stationsSlug, showId string) (*feeds.Feed, error) {
	eps, err := u.GetEpisodes(ctx, stationsSlug, showId)
	if err != nil {
		return nil, err
	}

	show, err := u.GetShow(ctx, stationsSlug, showId)
	if err != nil {
		return nil, err
	}

	return feeds2.ToFeed(u.hc, show, eps, eps[0].Description)
}

func (u *useCase) GetAllShowsFeed(ctx context.Context, stationsSlug string) (*feeds.Feed, error) {
	st, err := u.GetStation(ctx, stationsSlug)
	if err != nil {
		return nil, err
	}

	shows, err := u.GetShows(ctx, stationsSlug)
	if err != nil {
		return nil, err
	}

	episodes := make([]*models.Episode, 0)
	epsChan := make(chan *models.Episode)
	nshows := int32(len(shows))

	// load eps for all shows concurrently
	var g errgroup.Group
	for _, show := range shows {
		show := show
		g.Go(func() error {
			// last one out closes shop
			defer func() {
				if atomic.AddInt32(&nshows, -1) == 0 {
					close(epsChan)
				}
			}()

			// fetch episodes and add to channel
			eps, err := u.GetEpisodes(ctx, stationsSlug, show.Id)
			if err != nil {
				return err
			}
			for _, ep := range eps {
				epsChan <- ep
			}
			return nil
		})
	}
	g.Go(func() error {
		for ep := range epsChan {
			episodes = append(episodes, ep)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return feeds2.ToFeed(u.hc, &models.Show{Name: st.Name, ImageUrl: st.ImageUrl}, episodes, st.Tagline)
}

func NewUseCase(gp globalplayer.GlobalPlayer, hc resty.HttpClient) UseCase {
	return &useCase{
		gp: gp,
		hc: hc,
	}
}

func NewHttpClient() resty.HttpClient {
	return http.DefaultClient
}

var ProviderSet = wire.NewSet(NewUseCase, NewHttpClient)
