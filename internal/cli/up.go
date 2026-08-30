package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/JurisDab/devtool/internal/config"
	"github.com/JurisDab/devtool/internal/orchestrator"
)

var serviceColors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgYellow),
	color.New(color.FgGreen),
	color.New(color.FgMagenta),
	color.New(color.FgBlue),
}

func newUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start every service defined in the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			colorFor := assignColors(cfg)

			lines, err := orchestrator.Up(cmd.Context(), cfg)
			if err != nil {
				return err
			}

			for line := range lines {
				c := colorFor[line.Service]
				prefix := c.Sprintf("[%s]", line.Service)

				switch {
				case line.IsStatus && line.IsErr:
					fmt.Printf("%s %s\n", prefix, color.New(color.FgRed, color.Bold).Sprint("✗ "+line.Text))
				case line.IsStatus:
					fmt.Printf("%s %s\n", prefix, color.New(color.FgGreen, color.Bold).Sprint("✓ "+line.Text))
				default:
					fmt.Printf("%s %s\n", prefix, line.Text)
				}
			}

			return nil
		},
	}
}

func assignColors(cfg *config.Config) map[string]*color.Color {
	m := make(map[string]*color.Color, len(cfg.Services))
	for i, svc := range cfg.Services {
		m[svc.Name] = serviceColors[i%len(serviceColors)]
	}
	return m
}
