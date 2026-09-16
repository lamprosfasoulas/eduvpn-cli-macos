package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eduvpn-cli",
	Short: "Headless eduVPN client for macOS",
	Long: `eduvpn-cli connects to an eduVPN or Let's Connect! server and manages
the resulting WireGuard tunnel via wg-quick, without needing the GUI client.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
