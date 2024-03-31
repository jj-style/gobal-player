package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) StationsUI(c *gin.Context) {
	stations, err := s.uc.GetStations(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "stations.tmpl", gin.H{
		"data": stations,
		"breadcrumbs": []struct {
			Name   string
			Url    string
			Active bool
		}{
			{Name: "Stations", Url: c.Request.RequestURI, Active: true},
		},
	})
}

func (s *Service) ShowsUI(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
	}
	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station, err := s.uc.GetStation(c, req.Slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shows, err := s.uc.GetShows(c, req.Slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "shows.tmpl", gin.H{
		"ctx":     c,
		"data":    shows,
		"station": station,
		"breadcrumbs": []struct {
			Name   string
			Url    string
			Active bool
		}{
			{Name: "Stations", Url: "/stations"},
			{Name: station.Name, Url: c.Request.RequestURI, Active: true},
		},
	})
}

func (s *Service) EpisodesUI(c *gin.Context) {
	type request struct {
		Slug string `uri:"slug" binding:"required"`
		Id   string `uri:"id" binding:"required"`
	}
	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station, err := s.uc.GetStation(c, req.Slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	show, err := s.uc.GetShow(c, req.Slug, req.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	episodes, err := s.uc.GetEpisodes(c, req.Slug, req.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "episodes.tmpl", gin.H{
		"ctx":  c,
		"data": episodes,
		"show": show,
		"breadcrumbs": []struct {
			Name   string
			Url    string
			Active bool
		}{
			{Name: "Stations", Url: "/stations"},
			{Name: station.Name, Url: fmt.Sprintf("/shows/%s", station.Slug)},
			{Name: show.Name, Url: fmt.Sprintf("/episodes/%s/%s", station.Slug, show.Id), Active: true},
		},
	})
}
