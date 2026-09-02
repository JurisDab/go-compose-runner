package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/JurisDab/go-compose-runner/internal/config"
	"github.com/JurisDab/go-compose-runner/internal/doctor"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the local dev environment for common problems",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			checks := doctor.Run(cfg)

			failed := 0
			for _, c := range checks {
				if c.OK {
					fmt.Printf("%s %s\n", color.New(color.FgGreen).Sprint("✓"), c.Name)
					continue
				}
				failed++
				fmt.Printf("%s %s %s\n", color.New(color.FgRed).Sprint("✗"), c.Name, color.New(color.FgHiBlack).Sprintf("(%s)", c.Detail))
			}

			if failed > 0 {
				return fmt.Errorf("%d check(s) failed", failed)
			}
			fmt.Println(color.New(color.FgGreen, color.Bold).Sprint("\nAll checks passed."))
			return nil
		},
	}
}
