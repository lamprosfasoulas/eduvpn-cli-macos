package cmd

import (
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the disconnect command
var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect the active eduVPN tunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.Disconnect()
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
