package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App    App
		ArtNet ArtNet
	}

	App struct {
		Port     string `env:"HTTP_PORT" envDefault:"8080"`
		LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	}

	ArtNet struct {
		LogLevel string `env:"ARTNET_LOG_LEVEL" envDefault:"info"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
