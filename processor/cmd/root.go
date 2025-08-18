package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	ErrNonFuncCMD = errors.New("non functional command")
)

func Execute() {
	rootCMD := createRootCMD()

	rootCMD.AddCommand(createMigrateCMD())
	rootCMD.AddCommand(createProcessCMD())

	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}

func createRootCMD() *cobra.Command {
	return &cobra.Command{
		Use: "processor",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, args)

			return ErrNonFuncCMD
		},
	}
}
