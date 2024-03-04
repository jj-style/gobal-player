package globalplayer

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/chromedp/chromedp"
)

const (
	globalPlayerQueryUrl = "https://www.globalplayer.com/live/lbc/uk/"
)

var (
	buildIdRe          = regexp.MustCompile(`"buildId":"(?P<buildId>.*?)"`)
	ErrBuildIdNotFound = errors.New("could not find build id")
)

// Gets the build id for global player of errors.
func GetBuildId() (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var data string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(globalPlayerQueryUrl),
		chromedp.OuterHTML("html", &data, chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}

	buildId := buildIdRe.FindStringSubmatch(data)
	if len(buildId) != 2 {
		return "", ErrBuildIdNotFound
	}
	return buildId[1], nil
}

// Checks if the build id is still valid
func CheckBuildId(id string) bool {
	checkUrl, err := url.JoinPath(baseUrl, id, "index.json")
	if err != nil {
		return false
	}
	resp, err := http.Get(checkUrl)
	if err != nil {
		return false
	}

	if resp.StatusCode >= 400 {
		return false
	}
	return true
}
