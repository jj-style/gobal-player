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
	stations, err := s.uc.GetStations(c)
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

	shows, err := s.uc.GetShows(c, req.Slug)
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

	eps, err := s.uc.GetEpisodes(c, req.Slug, req.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"episodes": eps})
}

func (s *Service) GetEpisodesRss(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
		Id   string `uri:"id" binding:"required"`
	}

	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed, err := s.uc.GetEpisodesFeed(c, req.Slug, req.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.Header().Set("content-type", "application/xml")
	c.String(200, rss)
}

func (s *Service) GetAllShowsRss(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
	}

	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed, err := s.uc.GetAllShowsFeed(c, req.Slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.Header().Set("content-type", "application/xml")
	c.String(200, rss)
}

func NewService(uc globalplayer.UseCase) *Service {
	return &Service{uc: uc}
}
