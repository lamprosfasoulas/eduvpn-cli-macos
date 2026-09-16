package cmd

import (
	"fmt"

	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect <alias|institute|url>",
	Short: "Connect to a saved alias, an institute access server by name, or a self-hosted server with --custom",
	RunE: func(cmd *cobra.Command, args []string) error {
		custom, err := cmd.Flags().GetBool("custom")
		if err != nil {
			return err
		}
		if len(args) != 1 {
			return fmt.Errorf("server is needed")
		}
		server := args[0]

		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.Connect(server, custom)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().Bool("custom", false, "Treat <server> as a self-hosted base URL instead of searching institute access servers by name")
}
