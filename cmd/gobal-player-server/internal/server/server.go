package server

import (
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

func NewServer(service *service.Service) *Server {
	r := gin.Default()

	addRoutes(r, service)

	return &Server{
		Router: r,
	}
}

var GlobalPlayerProvider = wire.NewSet(NewCache, NewGlobalPlayer)

func NewCache(config *config.Config) resty.Cache[[]byte] {
	return resty.NewCache[[]byte](config.Cache.Ttl)
}

func NewGlobalPlayer(config *config.Config, cache resty.Cache[[]byte]) (globalplayer.GlobalPlayer, error) {
	buildId, err := globalplayer.GetBuildId(http.DefaultClient)
	if err != nil {
		return nil, err
	}
	return globalplayer.NewClient(http.DefaultClient, buildId, cache), nil
}
