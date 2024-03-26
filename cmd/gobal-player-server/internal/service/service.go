package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
)

type GlobalPlayerService interface {
	GetStations(c *gin.Context)
}

type service struct {
	uc globalplayer.UseCase
}

func (s *service) GetStations(c *gin.Context) {
	stations, err := s.uc.GetStations(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSONP(http.StatusOK, gin.H{"stations": stations})
}

func NewService(uc globalplayer.UseCase) GlobalPlayerService {
	return &service{uc: uc}
}
