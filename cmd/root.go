package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

var shell string

var rootCmd = &cobra.Command {
	Use: "devcontainer-shell",
	Short: "devcontainer shell helper",
	Run: func(cmd *cobra.Command, args []string) {
		if err := devcontainershell.DevcontainerExec(shell); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&shell, "shell", "s", "bash", "using shell program")

	rootCmd.AddCommand(execCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
