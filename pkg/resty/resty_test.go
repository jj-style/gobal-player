package resty_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/jj-style/gobal-player/pkg/resty/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_client_Post(t *testing.T) {
	var (
		baseUrl = "http://example.com/api"
	)
	type args struct {
		uri  string
		body any
	}
	tests := []struct {
		name     string
		args     args
		setup    func(hc *mocks.MockHttpClient)
		want     []byte
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy json",
			args: args{
				uri: "todo",
				body: struct {
					Name string `json:"name"`
					Done bool   `json:"done"`
				}{
					Name: "take bins out",
					Done: false,
				},
			},
			setup: func(hc *mocks.MockHttpClient) {
				var resp = &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("hello")),
				}
				hc.EXPECT().Do(mock.Anything).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					contentType := req.Header.Get("Content-Type")
					if contentType != "application/json" {
						t.Error(fmt.Errorf("wanted json content type, got: %s", contentType))
						t.FailNow()
					}
					return resp, nil
				})
			},
			want:     []byte("hello"),
			checkErr: assert.NoError,
		},
		{
			name: "happy form",
			args: args{
				uri:  "todo",
				body: url.Values{"q": []string{"searchTerm"}, "sortBy": []string{"desc"}},
			},
			setup: func(hc *mocks.MockHttpClient) {
				var resp = &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("hello")),
				}
				hc.EXPECT().Do(mock.Anything).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					contentType := req.Header.Get("Content-Type")
					if contentType != "application/x-www-form-urlencoded" {
						t.Error(fmt.Errorf("wanted form content type, got: %s", contentType))
						t.FailNow()
					}
					return resp, nil
				})
			},
			want:     []byte("hello"),
			checkErr: assert.NoError,
		},
		{
			name: "unhappy",
			args: args{
				uri:  "todo",
				body: nil,
			},
			setup: func(hc *mocks.MockHttpClient) {
				hc.EXPECT().Do(mock.Anything).Return(nil, errors.New("boom"))
			},
			want:     nil,
			checkErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := mocks.NewMockHttpClient(t)
			if tt.setup != nil {
				tt.setup(hc)
			}

			c := resty.NewClient(resty.WithHttpClient(hc), resty.WithBaseUrl(baseUrl))

			got, err := c.Post(tt.args.uri, tt.args.body)
			tt.checkErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_client_Get(t *testing.T) {
	var (
		baseUrl = "http://example.com/api"
	)
	type args struct {
		uri string
	}
	tests := []struct {
		name     string
		args     args
		setup    func(hc *mocks.MockHttpClient)
		want     []byte
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			args: args{
				uri: "test/1",
			},
			setup: func(hc *mocks.MockHttpClient) {
				var resp = &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("hello")),
				}
				hc.EXPECT().Do(mock.Anything).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						t.Errorf("expected method: %s, got: %s", http.MethodGet, req.Method)
					}
					wantUrl := "http://example.com/api/test/1"
					if req.URL.String() != wantUrl {
						t.Errorf("expexted url: %s, got: %s", wantUrl, req.URL.String())
					}
					return resp, nil
				})
			},
			want:     []byte("hello"),
			checkErr: assert.NoError,
		},
		{
			name: "unhappy",
			args: args{
				uri: "test/1",
			},
			setup: func(hc *mocks.MockHttpClient) {
				hc.EXPECT().Do(mock.Anything).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						t.Errorf("expected method: %s, got: %s", http.MethodGet, req.Method)
					}
					wantUrl := "http://example.com/api/test/1"
					if req.URL.String() != wantUrl {
						t.Errorf("expexted url: %s, got: %s", wantUrl, req.URL.String())
					}
					return nil, errors.New("boom")
				})
			},
			want:     nil,
			checkErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := mocks.NewMockHttpClient(t)
			if tt.setup != nil {
				tt.setup(hc)
			}

			c := resty.NewClient(resty.WithHttpClient(hc), resty.WithBaseUrl(baseUrl))

			got, err := c.Get(tt.args.uri)
			tt.checkErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientCacheSetCalledOnCacheMiss(t *testing.T) {
	hc := mocks.NewMockHttpClient(t)
	cache := mocks.NewMockCache[[]byte](t)
	var resp = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("hello")),
	}

	// cache miss
	cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
	// mock request/response
	hc.EXPECT().Do(mock.Anything).Return(resp, nil)
	// cache stored
	cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("string"), []byte("hello")).Return(nil)

	c := resty.NewClient(resty.WithHttpClient(hc), resty.WithCache(cache))

	c.Get("http://example.com")
}

func TestClientCacheGetHit(t *testing.T) {
	hc := mocks.NewMockHttpClient(t)
	cache := mocks.NewMockCache[[]byte](t)

	// cache hit
	cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("string")).Return([]byte("hello"), nil)

	c := resty.NewClient(resty.WithHttpClient(hc), resty.WithCache(cache))

	resp, err := c.Get("http://example.com")

	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), resp)
}
