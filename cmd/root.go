package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "devcontainer-shell",
	Short: "devcontainer shell helper",
	Run: exec,
}

func init() {
	setupExecCmd(rootCmd)

	rootCmd.AddCommand(execCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
