package pubsub

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/mmcdole/gofeed"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
)

type PubSub interface {
	Run(chan<- error) error
}

type pubsub struct {
	apiUrl string
	gp     globalplayer.GlobalPlayer
	rd     *redis.Client
	nc     *nats.Conn
	cron   *cron.Cron
}

func NewPubSub(apiUrl string, gp globalplayer.GlobalPlayer, rd *redis.Client, nc *nats.Conn) PubSub {
	cron := cron.New()
	cron.Start()
	return &pubsub{apiUrl, gp, rd, nc, cron}
}

func (x *pubsub) Run(errc chan<- error) error {
	stations, err := x.gp.GetStations()
	if err != nil {
		return err
	}

	// process all stations on a schedule looking for new episodes
	for _, st := range stations {
		st := st
		jobId, _ := x.cron.AddFunc("@every 10m", func() {
			go func() { x.processFeed(context.Background(), st, errc) }()
		})
		go x.cron.Entry(jobId).WrappedJob.Run()
	}
	return nil
}

func (x *pubsub) processFeed(ctx context.Context, station *models.Station, errc chan<- error) {
	feedUrl, _ := url.JoinPath(x.apiUrl, "shows", station.Slug, "rss")
	var feed *gofeed.Feed

	// try a few times to fetch the station feed
	if _, _, err := lo.AttemptWithDelay(3, time.Millisecond*500, func(index int, duration time.Duration) error {
		var err error
		feed, err = gofeed.NewParser().ParseURL(feedUrl)
		if err != nil {
			return fmt.Errorf("getting feed for %s: %w", station.Slug, err)
		}
		return nil
	}); err != nil {
		errc <- fmt.Errorf("getting feed for %s: %w", station.Slug, err)
		return
	}

	// for each item in the feed, index it in redis and publish it to nats
	for _, item := range feed.Items {
		if seen := x.rd.SIsMember(ctx, "items", item.GUID).Val(); !seen {
			if err := x.rd.SAdd(ctx, "items", item.GUID); err != nil {
				if err := x.nc.Publish(fmt.Sprintf("episodes.%s", station.Slug), []byte(item.Enclosures[0].URL)); err != nil {
					errc <- err
					continue
				}
			}
		}
	}
}
