package config

import (
	"context"

	cbConfig "github.com/go-coldbrew/core/config"
	"github.com/go-coldbrew/log"
	"github.com/kelseyhightower/envconfig"
)

var defaultConfig Config

type Config struct {
	cbConfig.Config
	PanicOnConfigError bool   `envconfig:"PANIC_ON_CONFIG_ERROR" default:"true"`
	DatabaseURL        string `envconfig:"DATABASE_URL" required:"true"`
	Prefix string `envconfig:"PREFIX" default:"got"`
}

func init() {
	err := envconfig.Process("", &defaultConfig)
	if err != nil {
		if defaultConfig.PanicOnConfigError {
			panic(err)
		} else {
			log.Error(context.Background(), "msg", "error while loading config", "err", err)
		}
	}
}

func Get() Config {
	return defaultConfig
}

func GetColdBrewConfig() cbConfig.Config {
	return defaultConfig.Config
}
