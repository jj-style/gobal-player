package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
)

func addRoutes(r *gin.Engine, service *service.Service) {
	api := r.Group("/api")
	api.GET("/stations", service.GetStations)
	api.GET("/shows/:slug", service.GetShows)
	api.GET("/episodes/:slug/:id", service.GetEpisodes)
	api.GET("/shows/:slug/rss", service.GetAllShowsRss)
	api.GET("/episodes/:slug/:id/rss", service.GetEpisodesRss)

	r.GET("/stations", loadTemplates(r), service.StationsUI)
	r.GET("/shows/:slug", loadTemplates(r), service.ShowsUI)
	r.GET("/episodes/:slug/:id", loadTemplates(r), service.EpisodesUI)
}
