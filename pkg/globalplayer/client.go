package globalplayer

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models/nextjs"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
)

const (
	loginUrl = "https://gigya.globalplayer.com/accounts.login"
	baseUrl  = "https://www.globalplayer.com/_next/data"
)

// GlobalPlayer is an interface to some of global players APIs
type GlobalPlayer interface {
	GetStations() ([]*models.Station, error)
	GetShows(stationSlug string) ([]*models.Show, error)
	GetEpisodes(stationSlug, showId string) ([]*models.Episode, error)
}

type gpClient struct {
	rc   resty.Client
	cron *cron.Cron
}

func NewClient(hc *http.Client, cache resty.Cache[[]byte], updateDuration string) (GlobalPlayer, func(), error) {

	newRestClient := func() (resty.Client, error) {
		buildId, err := GetBuildId(hc)
		if err != nil {
			return nil, err
		}
		baseUrlWithApiKey, _ := url.JoinPath(baseUrl, buildId)
		rc := resty.NewClient(
			resty.WithBaseUrl(baseUrlWithApiKey),
			resty.WithHttpClient(hc),
			resty.WithCache(cache),
		)
		return rc, nil
	}

	rc, err := newRestClient()
	if err != nil {
		return nil, func() {}, err
	}

	cron := cron.New()
	client := &gpClient{rc: rc, cron: cron}

	if updateDuration != "" {
		_, err = cron.AddFunc("@every 1m", func() {
			if rc, err = newRestClient(); err == nil {
				client.rc = rc
			}
		})
		if err != nil {
			return nil, func() {}, err
		}
	}

	cron.Start()

	return client, func() { cron.Stop() }, nil
}

func (c *gpClient) GetStations() ([]*models.Station, error) {
	resp, err := resty.Get[nextjs.StationsPageResponse](c.rc, "index.json")
	if err != nil {
		return nil, err
	}

	brands := resp.PageProps.Feature.Blocks[0].Brands
	return lo.Map(brands, func(item nextjs.StationBrand, _ int) *models.Station {
		return models.StationFromApiModel(&item)
	}), nil
}

func (c *gpClient) GetShows(stationSlug string) ([]*models.Show, error) {
	u, _ := url.JoinPath("catchup", stationSlug, "uk.json")
	resp, err := resty.Get[nextjs.CatchupResponse](c.rc, u)
	if err != nil {
		return nil, err
	}

	return lo.Map(resp.PageProps.CatchupInfo, func(item nextjs.CatchupInfo, _ int) *models.Show {
		return models.ShowFromApiModel(&item)
	}), nil
}

func (c *gpClient) GetEpisodes(stationSlug, showId string) ([]*models.Episode, error) {
	u, _ := url.JoinPath("catchup", stationSlug, "uk", fmt.Sprintf("%s.json", showId))
	resp, err := resty.Get[nextjs.CatchupShowResponse](c.rc, u)
	if err != nil {
		return nil, err
	}

	return lo.Map(resp.PageProps.CatchupInfo.Episodes, func(item nextjs.Episode, _ int) *models.Episode {
		return models.EpisodeFromApiModel(&item)
	}), nil
}

// Login logs in through the the global player API returning the authorisation response, or errors.
func Login(email, password, apiKey string) (nextjs.LoginResponse, error) {
	return resty.Post[url.Values, nextjs.LoginResponse](
		resty.NewClient(),
		loginUrl,
		url.Values{
			"LoginId":  []string{email},
			"Password": []string{password},
			"APIKey":   []string{apiKey},
		},
	)
}
