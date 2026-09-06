package cli

import (
	"github.com/spf13/cobra"
)

var configPath string

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "devtool",
		Short: "devtool automates local multi-service dev environments",
		Long: "devtool replaces ad-hoc shell scripts for local development: " +
			"one command to start every service, tail their logs together, " +
			"and check your environment for common problems.",
	}

	root.PersistentFlags().StringVarP(&configPath, "config", "c", ".devtool.yaml", "path to devtool config file")
	root.SilenceUsage = true
	root.SilenceErrors = true

	root.AddCommand(newUpCmd())
	root.AddCommand(newDownCmd())
	root.AddCommand(newLogsCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newNewCmd())

	return root
}
