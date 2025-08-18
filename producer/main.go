package main

import (
	"github.com/EWK20/event-processor/producer/config"
	"github.com/EWK20/event-processor/producer/producer"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	producer, err := producer.New(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create producer")
	}

	producer.Run()
}
