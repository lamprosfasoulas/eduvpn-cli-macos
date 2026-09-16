package cmd

import (
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the eduVPN tunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.Status()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
