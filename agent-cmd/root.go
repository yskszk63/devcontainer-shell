package agentcmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "devcontainer-shell-agent",
	SilenceErrors: true,
	SilenceUsage: true,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(watchListensCmd)
	rootCmd.AddCommand(installCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
