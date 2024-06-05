package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jj-style/gobal-player/cmd/podpub/internal/pubsub"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var (
	serverUrl string
	redisUrl  string
	natsUrl   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "periodically fetch all rss feeds and send notifications for new episodes",
	RunE:  run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&serverUrl, "server", "s", "", "base url of gobal-player server api")
	rootCmd.Flags().StringVarP(&redisUrl, "redis", "r", "", "redis url")
	rootCmd.Flags().StringVarP(&natsUrl, "nats", "n", "", "nats url")

	rootCmd.MarkFlagRequired("server")
	rootCmd.MarkFlagRequired("redis")
	rootCmd.MarkFlagRequired("nats")
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	gp, cleanup, err := globalplayer.NewClient(http.DefaultClient, resty.NewCache[[]byte](0), "1h")
	if err != nil {
		return err
	}
	defer cleanup()

	rd := redis.NewClient(&redis.Options{Addr: redisUrl})

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}
	defer func() {
		nc.Drain()
	}()

	pubsub := pubsub.NewPubSub(serverUrl, gp, rd, nc)

	errc := make(chan error)
	if err := pubsub.Run(errc); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		close(errc)
	}()

	for err := range errc {
		log.Println(err)
	}

	return nil
}
