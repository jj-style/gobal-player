package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/config"
)

var (
	confDir string
)

func init() {
	flag.StringVar(&confDir, "conf", "", "-conf /dir/with/config/in")
	flag.Parse()
}

func run(ctx context.Context, config *config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	server, cleanup, err := InitializeServer(config)
	if err != nil {
		return err
	}
	defer cleanup()

	go func() {
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		log.Printf("listening on %s\n", addr)
		if err := server.Router.Run(addr); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}()
	wg.Wait()
	return nil
}

func main() {
	conf, err := config.NewConfig(confDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v", err)
		os.Exit(1)
	}
	if err := run(context.Background(), conf); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
