package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	appName = "gobal-player-tui"
)

var (
	homedir    = MustHomeDir()
	UserAppDir = path.Join(homedir, "."+appName)
	C          Config
)

type Config struct {
	BuildId  string `mapstructure:"buildId"`
	Insecure bool   `mapstructure:"insecure"`
}

func NewConfig() *Config {
	return &Config{Insecure: false}
}

// Initailises viper and loads the configuration
func InitConfig() {
	if err := os.Mkdir(UserAppDir, os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}

	viper.SetConfigName("config")                   // name of config file (without extension)
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path.Join("/etc", appName)) // path to look for the config file in
	viper.AddConfigPath(UserAppDir)                 // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory
	viper.AutomaticEnv()                            // automatically merge in environment variables
	viper.SetEnvPrefix("gp")                        // only consider environment variables starting "GP_"

	viper.SetDefault("insecure", "false")

	ReadInConfig()
}

// ReadInConfig reads the config based on what is initialised and
// parses it into the global config struct
func ReadInConfig() {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create default
			C = *NewConfig()
		} else {
			// Config file was found but another error was produced
			log.Fatal(fmt.Errorf("fatal error config file: %w", err))
		}
	} else {
		// deserialise config into global struct
		err = viper.Unmarshal(&C)
		if err != nil {
			log.Fatal(fmt.Errorf("unable to decode into struct, %v", err))
		}
		fmt.Printf("%+v\n", C)
		fmt.Printf("%+v\n", viper.AllSettings())
	}
}
