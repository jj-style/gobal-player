package feeds

import (
	"fmt"

	"github.com/gorilla/feeds"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/samber/lo"
)

// Given a show and it's episodes, return a feed to consume the show.
func ToFeed(show *models.Show, episodes []*models.Episode, description string) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       show.Name,
		Image:       &feeds.Image{Url: show.ImageUrl},
		Updated:     lo.MaxBy(episodes, func(a, b *models.Episode) bool { return b.Aired.After(a.Aired) }).Aired,
		Description: description,
	}
	feed.Items = lo.Map(episodes, func(item *models.Episode, _ int) *feeds.Item {
		return &feeds.Item{
			Title:       fmt.Sprintf("%s: %s", item.Name, item.Aired.Format("Monday 02 January 2006")),
			Link:        &feeds.Link{Href: item.StreamUrl},
			Id:          item.Id,
			Created:     item.Aired,
			Description: fmt.Sprintf("%s<br/><br/>Available until %s.", item.Description, item.Until.Format("Monday 02 January 2006 15:04:05")),
			Enclosure:   &feeds.Enclosure{Url: item.StreamUrl, Type: "audio/mpeg", Length: "1"},
		}
	})

	return feed
}
