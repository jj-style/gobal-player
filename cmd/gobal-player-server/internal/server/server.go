package server

import (
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/config"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/resty"
)

type Server struct {
	Router *gin.Engine
}

func NewServer(service service.GlobalPlayerService) *Server {
	r := gin.Default()

	addRoutes(r, service)

	return &Server{
		Router: r,
	}
}

var GlobalPlayerProvider = wire.NewSet(NewHttpClient, NewCache, NewGlobalPlayer)

func NewCache(config *config.Config) resty.Cache[[]byte] {
	return resty.NewCache[[]byte](config.Cache.Ttl)
}

func NewGlobalPlayer(config *config.Config, hc *http.Client, cache resty.Cache[[]byte]) (globalplayer.GlobalPlayer, error) {
	buildId, err := globalplayer.GetBuildId(hc)
	if err != nil {
		return nil, err
	}
	return globalplayer.NewClient(hc, buildId, cache), nil
}

func NewHttpClient(config *config.Config) *http.Client {
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure}}}
}
