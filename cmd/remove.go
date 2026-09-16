package cmd

import (
	"fmt"

	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/handlers"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove a saved server alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("alias is needed")
		}

		client, err := handlers.NewClient()
		if err != nil {
			return err
		}
		defer client.Cancel()
		return client.RemoveAlias(args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
