package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host  string      `mapstructure:"host"`
	Port  int         `mapstructure:"port"`
	Cache CacheConfig `mapstructure:"cache"`
}

type CacheConfig struct {
	Ttl time.Duration `mapstructure:"ttl"`
}

func NewConfig(dir string) (*Config, error) {
	// setup some defaults
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8080")
	viper.SetDefault("cache.ttl", "60s")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("gp")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}
