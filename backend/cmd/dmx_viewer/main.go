package main

import (
	"log"

	"github.com/nasshu2916/dmx_viewer/internal/app"
	"github.com/nasshu2916/dmx_viewer/internal/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	app.Run(config)
}
