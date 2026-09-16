package cmd

import (
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "List saved server aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.ListAliases()
	},
}

func init() {
	rootCmd.AddCommand(serversCmd)
}
