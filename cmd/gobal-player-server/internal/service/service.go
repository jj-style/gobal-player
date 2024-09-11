package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/feeds"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/samber/lo"
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

func (s *Service) GetLive(c *gin.Context) {
	live, err := s.uc.GetLive(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"live": live})
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

	liveStations, err := s.uc.GetLive(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get the live shows, and if there is a corresponding live feed of the same station
	// then add a podcastin2.0 live item tag to the feed
	if live, found := lo.Find(liveStations, func(item *models.LiveStation) bool {
		return item.Slug == req.Slug
	}); found {
		feed.Podcasting2LiveItem = &feeds.Podcasting2LiveItem{
			Status: "live",
			Start:  time.Now().Add(-time.Hour).Format(time.RFC3339),
			End:    time.Now().Add(time.Hour).Format(time.RFC3339),
			RssItem: &feeds.RssItem{
				Enclosure:   &feeds.RssEnclosure{Url: live.StreamUrl, Type: "audio/mpeg", Length: "312"},
				Title:       fmt.Sprintf("%s Live!", live.Name),
				Description: live.Tagline,
				Link:        live.StreamUrl,
				Guid:        &feeds.RssGuid{Id: live.StreamUrl, IsPermaLink: "true"},
				Podcasting2Item: &feeds.Podcasting2Item{
					ContentLink: &feeds.Podcasting2ContentLink{
						Href: live.StreamUrl,
						Text: "Listen Live!",
					},
					AlternateEnclosure: &feeds.Podcasting2AlternateEnclosure{
						Type:    "audio/mpeg",
						Length:  "312",
						Default: true,
						Source: feeds.Podcasting2Source{
							Uri: live.StreamUrl,
						},
					},
				},
			},
		}
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
