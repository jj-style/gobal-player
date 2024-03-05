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
)

func init() {
	config.InitConfig()
}

func main() {
	cleanup1 := initLogger()
	defer cleanup1()

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: config.C.Insecure}}}
	checkOrRegenConfig(httpClient)
	fmt.Printf("insecure %+v\n", config.C.Insecure)
	gp := globalplayer.NewClient(httpClient, config.C.BuildId)

	player, cleanup, err := audioplayer.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app := NewApp(gp, player)
	if err := app.Run(); err != nil {
		log.Fatal(err)
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

// checks the state of the config
// and regenerates a new buildId if it is invalid
func checkOrRegenConfig(hc *http.Client) {
	if ok := globalplayer.CheckBuildId(hc, config.C.BuildId); !ok {
		newBuildId, err := globalplayer.GetBuildId(hc)
		if err != nil {
			log.Fatal(err)
		}
		viper.Set("buildId", newBuildId)
		if err := viper.WriteConfigAs(path.Join(config.UserAppDir, "config.yaml")); err != nil {
			log.Fatal(err)
		}
		config.ReadInConfig()
	}
}
