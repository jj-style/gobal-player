package feeds

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dannav/hhmmss"
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
		Subtitle:    description,
		ITunes: &feeds.ITunesFeed{
			Title:    show.Name,
			Image:    &feeds.ITunesImage{Href: show.ImageUrl},
			Type:     feeds.ITunesFeedTypeEpisodic,
			Explicit: false,
		},
	}

	feedItems := make([]*feeds.Item, len(episodes))
	var g errgroup.Group
	for idx, item := range episodes {
		idx := idx
		item := item
		g.Go(func() error {
			duration, err := hhmmss.Parse(item.Duration)
			if err != nil {
				return fmt.Errorf("parsing episode duration '%s': %v", item.Duration, err)
			}

			lengthChan := lo.Async(func() int64 {
				headReq, _ := http.NewRequest(http.MethodHead, item.StreamUrl, nil)
				streamHead, err := hc.Do(headReq)
				if err != nil {
					return 1
				}
				return streamHead.ContentLength
			})

			var contentLength int64 = 1
			select {
			case <-time.After(time.Second * 5):
			case contentLength = <-lengthChan:
			}

			feedItems[idx] = &feeds.Item{
				Title:       fmt.Sprintf("%s: %s", item.Name, item.Aired.Format("Monday 02 January 2006")),
				Link:        &feeds.Link{Href: item.StreamUrl},
				Id:          item.Id,
				Created:     item.Aired,
				Description: fmt.Sprintf("%s<br/><br/>Available until %s.", item.Description, item.Until.Format("Monday 02 January 2006 15:04:05")),
				Enclosure:   &feeds.Enclosure{Url: item.StreamUrl, Type: "audio/mpeg", Length: fmt.Sprint(contentLength)},
				ITunes: &feeds.ITunesItem{
					Duration:    fmt.Sprint(duration),
					EpisodeType: feeds.ITunesEpisodeTypeFull,
				},
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
