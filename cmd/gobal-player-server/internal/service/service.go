package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
)

type Service struct {
	uc globalplayer.UseCase
}

func (s *Service) GetStations(c *gin.Context) {
	stations, err := s.uc.GetStations(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"stations": stations})
}

func (s *Service) GetShows(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
	}

	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shows, err := s.uc.GetShows(c.Request.Context(), req.Slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"shows": shows})
}

func (s *Service) GetEpisodes(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
		Id   string `uri:"id" binding:"required"`
	}

	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eps, err := s.uc.GetEpisodes(c.Request.Context(), req.Slug, req.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"episodes": eps})
}

func NewService(uc globalplayer.UseCase) *Service {
	return &Service{uc: uc}
}
