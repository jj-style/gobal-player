package resty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Inteface of http.Client
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client interface {
	Post(uri string, body any) ([]byte, error)
	Get(uri string) ([]byte, error)
}

type client struct {
	http    HttpClient
	baseUrl *url.URL
	cache   Cache[[]byte]
}

type ClientOption func(*client)

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *client) {
		u, err := url.Parse(baseUrl)
		if err != nil {
			log.Fatal(err)
		}
		c.baseUrl = u
	}
}

func WithCache(ca Cache[[]byte]) ClientOption {
	return func(c *client) {
		c.cache = ca
	}
}

func WithHttpClient(hc HttpClient) ClientOption {
	return func(c *client) {
		c.http = hc
	}
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		http:  &http.Client{},
		cache: &nilCache[[]byte]{},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *client) doRestRequest(req *http.Request) (*http.Response, error) {
	log.WithFields(log.Fields{"method": req.Method, "uri": req.URL}).Debug("making request")
	resp, err := c.http.Do(req)
	if err != nil {
		log.Error(err.Error())
	}
	return resp, err
}

func (c *client) do(req *http.Request) ([]byte, error) {
	key := cacheKey(req)
	if hit, _ := c.cache.Get(context.TODO(), key); hit != nil {
		return hit, nil
	}

	resp, err := c.doRestRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = c.cache.Set(context.TODO(), key, out)

	if resp.StatusCode > 300 {
		log.WithFields(log.Fields{"code": resp.StatusCode, "url": req.URL}).Error("bad status code")
	}
	return out, nil
}

func (c *client) Post(uri string, body any) ([]byte, error) {
	_url, err := url.Parse(uri)
	if err == nil {
		if _url.Scheme == "" {
			_url = c.baseUrl.JoinPath(uri)
		}
	} else {
		return nil, err
	}

	var bodyData io.Reader
	var contentType string
	switch dataIn := body.(type) {
	case url.Values:
		bodyData = strings.NewReader(dataIn.Encode())
		contentType = "application/x-www-form-urlencoded"
	default:
		bin, err := json.Marshal(dataIn)
		if err != nil {
			return nil, err
		}
		bodyData = bytes.NewBuffer(bin)
		contentType = "application/json"
	}

	req, err := http.NewRequest(http.MethodPost, _url.String(), bodyData)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.do(req)
}

func (c *client) Get(uri string) ([]byte, error) {
	_url := uri
	if c.baseUrl != nil {
		_url = c.baseUrl.JoinPath(_url).String()
	}
	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func cacheKey(req *http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL)
}
