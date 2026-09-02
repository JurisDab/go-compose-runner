package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/JurisDab/go-compose-runner/internal/scaffold"
)

func newNewCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "new <template> <name>",
		Short: "Scaffold a new service from an embedded template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName, name := args[0], args[1]

			if err := scaffold.Generate(templateName, name, scaffold.Data{Name: name, Port: port}); err != nil {
				return err
			}

			fmt.Printf("Created %s/ from template %q\n", name, templateName)
			fmt.Printf("Add it to your .devtool.yaml, e.g.:\n\n")
			fmt.Printf("  - name: %s\n    workDir: ./%s\n    port: %d\n", name, name, port)
			return nil
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "port the new service should listen on")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := scaffold.List()
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(names, "\n"))
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	return cmd
}
