package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
)

func addRoutes(r *gin.Engine, service service.GlobalPlayerService) {
	r.GET("/stations", service.GetStations)
}
