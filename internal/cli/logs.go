package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/JurisDab/devtool/internal/config"
	"github.com/JurisDab/devtool/internal/orchestrator"
)

func newLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs [service...]",
		Short: "Tail logs for one or all services, multiplexed into one stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			colorFor := assignColors(cfg)

			lines, err := orchestrator.Logs(cmd.Context(), cfg, args)
			if err != nil {
				return err
			}

			for line := range lines {
				c := colorFor[line.Service]
				if c == nil {
					c = color.New(color.FgWhite)
				}
				prefix := c.Sprintf("[%s]", line.Service)
				fmt.Printf("%s %s\n", prefix, line.Text)
			}

			return nil
		},
	}
}
