package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
)

func addRoutes(r *gin.Engine, service *service.Service) {
	r.GET("/stations", service.GetStations)
	r.GET("/shows/:slug", service.GetShows)
	r.GET("/episodes/:slug/:id", service.GetEpisodes)
	r.GET("/shows/:slug/rss", service.GetAllShowsRss)
	r.GET("/episodes/:slug/:id/rss", service.GetEpisodesRss)
}
