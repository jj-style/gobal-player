package globalplayer

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const (
	globalPlayerQueryUrl = "https://www.globalplayer.com/live/lbc/uk/"
)

var (
	buildIdRe          = regexp.MustCompile(`"buildId":"(?P<buildId>.*?)"`)
	ErrBuildIdNotFound = errors.New("could not find build id")
)

// Gets the build id for global player of errors.
func GetBuildId(client *http.Client) (string, error) {
	resp, err := client.Get(globalPlayerQueryUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	buildId := buildIdRe.FindStringSubmatch(string(b))
	if len(buildId) != 2 {
		return "", ErrBuildIdNotFound
	}
	return buildId[1], nil
}

// Checks if the build id is still valid
func CheckBuildId(client *http.Client, id string) bool {
	checkUrl, err := url.JoinPath(baseUrl, id, "index.json")
	if err != nil {
		return false
	}
	resp, err := client.Get(checkUrl)
	if err != nil {
		return false
	}

	if resp.StatusCode >= 400 {
		return false
	}
	return true
}
