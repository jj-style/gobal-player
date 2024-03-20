package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	appName = "gobal-player-tui"
)

var (
	homedir    = MustHomeDir()
	UserAppDir = path.Join(homedir, "."+appName)
)

// Initailises viper and loads the configuration
func InitConfig() {
	if err := os.Mkdir(UserAppDir, os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}

	viper.SetDefault("insecure", "false")

	viper.SetConfigName("config")                   // name of config file (without extension)
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path.Join("/etc", appName)) // path to look for the config file in
	viper.AddConfigPath(UserAppDir)                 // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory
	viper.AutomaticEnv()                            // automatically merge in environment variables
	viper.SetEnvPrefix("gp")                        // only consider environment variables starting "GP_"
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("Config file changed: %s", e.Name)
		ReadInConfig()
	})
	viper.WatchConfig()

	ReadInConfig()
}

// ReadInConfig reads the config based on what is initialised and
// parses it into the global config struct
func ReadInConfig() {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create default
			if err = viper.WriteConfigAs(path.Join(UserAppDir, "config.yaml")); err != nil {
				log.Fatal(fmt.Errorf("writing default config file: %v", err))
			}
		} else {
			// Config file was found but another error was produced
			log.Fatal(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}
