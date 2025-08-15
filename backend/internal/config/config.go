package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App    App
		ArtNet ArtNet
		NTP    NTP
	}

	App struct {
		Port               string `env:"HTTP_PORT" envDefault:"8080"`
		LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
		HTTPTimeoutSeconds int    `env:"HTTP_TIMEOUT_SECONDS" envDefault:"30"`
	}

	ArtNet struct {
		LogLevel            string `env:"ARTNET_LOG_LEVEL" envDefault:"info"`
		ShortName           string `env:"ARTNET_SHORT_NAME" envDefault:"DMX Viewer"`
		LongName            string `env:"ARTNET_LONG_NAME" envDefault:"DMX Viewer Application"`
		PollIntervalSeconds int    `env:"ARTNET_POLL_INTERVAL_SECONDS" envDefault:"5"`
		ChannelBufferSize   int    `env:"ARTNET_CHANNEL_BUFFER_SIZE" envDefault:"1000"`
	}

	NTP struct {
		Enabled               bool   `env:"NTP_ENABLED" envDefault:"true"`
		Server                string `env:"NTP_SERVER" envDefault:"pool.ntp.org"`
		UpdateIntervalMinutes int    `env:"NTP_UPDATE_INTERVAL_MINUTES" envDefault:"360"`
		RetryCount            int    `env:"NTP_RETRY_COUNT" envDefault:"3"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
