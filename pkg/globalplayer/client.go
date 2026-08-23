package globalplayer

import (
	"fmt"
	"log"
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
	bffUrl   = "https://bff-web-guacamole.musicradio.com"
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

	bffClient resty.Client
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

	bffClient := resty.NewClient(
		resty.WithBaseUrl(bffUrl),
		resty.WithCache(cache),
	)

	client := &gpClient{rc: rc, cron: cron, bffClient: bffClient}

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

	catchupBlocks := resp.PageProps.CatchupInfo.Blocks
	listingsBlock, ok := lo.Find(catchupBlocks, func(item nextjs.CatchupBlock) bool { return item.Type == "Listing" })
	episodes := make([]*models.Episode, 0, len(catchupBlocks))
	if !ok {
		return episodes, nil
	}

	for _, listing := range listingsBlock.Items {
		playable, err := resty.Get[nextjs.BffPlayableResponse](c.bffClient, "/playables/"+listing.ID)
		if err != nil {
			log.Printf("error getting playable: %v", err)
			continue
		}

		var streamUrl string
		for _, playback := range playable.Playback {
			// TODO - get authentication for AdFree Url
			if !lo.Contains(playback.Flags, "AdFree") && playback.Url != "" {
				streamUrl = playback.Url
				break
			}
		}

		episodes = append(episodes, &models.Episode{
			Id:              playable.Id,
			Name:            playable.Title,
			Description:     listing.Description,
			ImageUrl:        listing.Image.Url,
			StreamUrl:       streamUrl,
			Duration:        listing.Content.Duration,
			DurationSeconds: listing.Content.DurationSeconds,
			Aired:           listing.Content.Published,
			Until:           listing.Content.Expiry,
			Availability:    listing.Content.Expiry.String(),
		})
	}
	return episodes, nil
}

// Login logs in through the the global player API returning the authorisation response, or errors.
func Login(email, password, apiKey string) (nextjs.LoginResponse, error) {
	return resty.Post[map[string]string, nextjs.LoginResponse](
		resty.NewClient(),
		loginUrl,
		map[string]string{
			"LoginId":           email,
			"password":          password,
			"APIKey":            apiKey,
			"includeUserInfo":   "true",
			"include":           "profile,data",
			"targetEnv":         "jssdk",
			"sdk":               "js_latest",
			"sessionExpiration": "0",
		},
	)
}
