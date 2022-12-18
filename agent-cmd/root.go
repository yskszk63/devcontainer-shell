package agentcmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "devcontainer-shell-agent",
	Run: func(c *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(watchListensCmd)
	rootCmd.AddCommand(installCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
