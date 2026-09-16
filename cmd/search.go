package cmd

import (
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "List institute access servers, optionally filtered by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.Search(query)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
