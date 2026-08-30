package cli

import (
	"github.com/spf13/cobra"

	"github.com/JurisDab/devtool/internal/config"
	"github.com/JurisDab/devtool/internal/orchestrator"
)

func newDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop every service defined in the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			return orchestrator.Down(cmd.Context(), cfg)
		},
	}
}
