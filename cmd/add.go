package cmd

import (
	"fmt"

	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <alias> <institute|url>",
	Short: "Save a server under an alias for later use with connect",
	RunE: func(cmd *cobra.Command, args []string) error {
		custom, err := cmd.Flags().GetBool("custom")
		if err != nil {
			return err
		}
		if len(args) != 2 {
			return fmt.Errorf("alias and server are needed")
		}
		alias, server := args[0], args[1]

		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.AddAlias(alias, server, custom)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().Bool("custom", false, "Treat <server> as a self-hosted base URL instead of searching institute access servers by name")
}
