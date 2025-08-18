package cmd

import (
	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/db"
	"github.com/EWK20/event-processor/processor/internal/processor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createProcessCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "process",
		Short: "Process events",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.New()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get config")
			}

			db, err := db.New(cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to database")
			}

			processor, err := processor.New(cfg.AWS, db)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to instantiate events processor")
			}

			processor.Run()
		},
	}
}
