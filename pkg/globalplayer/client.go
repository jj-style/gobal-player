package globalplayer

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/resty"
)

const (
	loginUrl = "https://gigya.globalplayer.com/accounts.login"
	baseUrl  = "https://www.globalplayer.com/_next/data"
)

// GlobalPlayer is an interface to some of global players APIs
type GlobalPlayer interface {
	GetStations() (models.StationsPageResponse, error)
	GetCatchup(slug string) (models.CatchupResponse, error)
	GetCatchupShows(slug, id string) (models.CatchupShowResponse, error)
	GetLive() (models.LiveResponse, error)
}

type gpClient struct {
	rc resty.Client
}

func NewClient(hc *http.Client, apiKey string) GlobalPlayer {
	baseUrlWithApiKey, _ := url.JoinPath(baseUrl, apiKey)
	restClient := resty.NewClient(
		resty.WithBaseUrl(baseUrlWithApiKey),
		resty.WithHttpClient(hc),
		resty.WithCache(resty.NewCache[[]byte](time.Hour)),
	)
	c := &gpClient{rc: restClient}
	return c
}

func (c *gpClient) GetStations() (models.StationsPageResponse, error) {
	return resty.Get[models.StationsPageResponse](c.rc, "index.json")
}

func (c *gpClient) GetCatchup(slug string) (models.CatchupResponse, error) {
	u, _ := url.JoinPath("catchup", slug, "uk.json")
	return resty.Get[models.CatchupResponse](c.rc, u)
}

func (c *gpClient) GetCatchupShows(slug, id string) (models.CatchupShowResponse, error) {
	u, _ := url.JoinPath("catchup", slug, "uk", fmt.Sprintf("%s.json", id))
	return resty.Get[models.CatchupShowResponse](c.rc, u)
}

func (c *gpClient) GetLive() (models.LiveResponse, error) {
	u, _ := url.JoinPath("live", "lbc", "uk.json")
	return resty.Get[models.LiveResponse](c.rc, u)
}

// Login logs in through the the global player API returning the authorisation response, or errors.
func Login(email, password, apiKey string) (models.LoginResponse, error) {
	return resty.Post[url.Values, models.LoginResponse](
		resty.NewClient(),
		loginUrl,
		url.Values{
			"LoginId":  []string{email},
			"Password": []string{password},
			"APIKey":   []string{apiKey},
		},
	)
}
