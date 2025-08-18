package cmd

import (
	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createMigrateCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.New()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get config")
			}

			db, err := db.New(cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to database")
			}

			if err := db.RunMigrations(); err != nil {
				log.Fatal().Err(err).Msg("failed to run database migrations")
			}
		},
	}
}
