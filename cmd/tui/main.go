package main

import (
	"fmt"
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
	if ok := globalplayer.CheckBuildId(config.C.BuildId); !ok {
		newBuildId, err := globalplayer.GetBuildId()
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

func main() {
	cleanup1 := initLogger()
	defer cleanup1()

	gp := globalplayer.NewClient(config.C.BuildId)

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
