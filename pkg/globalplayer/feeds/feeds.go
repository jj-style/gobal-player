package feeds

import (
	"fmt"
	"net/http"

	"github.com/gorilla/feeds"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

// Given a show and it's episodes, return a feed to consume the show.
func ToFeed(hc resty.HttpClient, show *models.Show, episodes []*models.Episode, description string) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title:       show.Name,
		Image:       &feeds.Image{Url: show.ImageUrl},
		Updated:     lo.MaxBy(episodes, func(a, b *models.Episode) bool { return b.Aired.After(a.Aired) }).Aired,
		Description: description,
	}

	feedItems := make([]*feeds.Item, len(episodes))
	var g errgroup.Group
	for idx, item := range episodes {
		idx := idx
		item := item
		g.Go(func() error {
			headReq, _ := http.NewRequest(http.MethodHead, item.StreamUrl, nil)
			streamHead, err := hc.Do(headReq)
			if err != nil {
				return fmt.Errorf("fetching episode HEAD: %v", err)
			}

			feedItems[idx] = &feeds.Item{
				Title:       fmt.Sprintf("%s: %s", item.Name, item.Aired.Format("Monday 02 January 2006")),
				Link:        &feeds.Link{Href: item.StreamUrl},
				Id:          item.Id,
				Created:     item.Aired,
				Description: fmt.Sprintf("%s<br/><br/>Available until %s.", item.Description, item.Until.Format("Monday 02 January 2006 15:04:05")),
				Enclosure:   &feeds.Enclosure{Url: item.StreamUrl, Type: "audio/mpeg", Length: fmt.Sprint(streamHead.ContentLength)},
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("constructing feed items: %v", err)
	}

	feed.Items = feedItems

	return feed, nil
}
