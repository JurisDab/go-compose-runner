// Package cli wires up devtool's Cobra command tree.
package cli

import (
	"github.com/spf13/cobra"
)

var configPath string

// NewRootCmd builds the devtool root command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "devtool",
		Short: "devtool automates local multi-service dev environments",
		Long: "devtool replaces ad-hoc shell scripts for local development: " +
			"one command to start every service, tail their logs together, " +
			"and check your environment for common problems.",
	}

	root.PersistentFlags().StringVarP(&configPath, "config", "c", ".devtool.yaml", "path to devtool config file")

	root.AddCommand(newUpCmd())
	root.AddCommand(newDownCmd())
	root.AddCommand(newLogsCmd())

	return root
}
