package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jj-style/gobal-player/cmd/gobal-player-tui/internal/config"
	"github.com/jj-style/gobal-player/pkg/audioplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/resty"
)

func init() {
	config.InitConfig()
}

func main() {
	cleanup1 := initLogger()
	defer cleanup1()

	httpClient, err := newHttpClient()
	if err != nil {
		log.Error(err)
		return
	}

	defer viper.WriteConfig()

	// don't expire cache in the TUI
	cache := resty.NewCache[[]byte](0)

	gp, cleanup2, err := globalplayer.NewClient(httpClient, cache, "@every 1h")
	if err != nil {
		log.Error(err)
		return
	}
	defer cleanup2()

	player, cleanup3, err := audioplayer.NewPlayer()
	if err != nil {
		log.Error(err)
		return
	}
	defer cleanup3()

	app := NewApp(gp, player, httpClient, cache)
	if err := app.Run(); err != nil {
		log.Error(err)
		return
	}
}

func initLogger() func() {
	f, err := os.OpenFile(path.Join(config.UserAppDir, "tui.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	l := log.StandardLogger()
	l.SetOutput(f)
	l.SetLevel(log.InfoLevel)
	return func() {
		// don't forget to close it
		f.Close()
	}
}

// creates a new *http.Client based on the config.
func newHttpClient() (*http.Client, error) {
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("insecure")}}}, nil
}
