package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

var execCmd = &cobra.Command {
	Use: "exec",
	Short: "execute shell on devcontainer",
	Run: func(cmd *cobra.Command, args []string) {
		if err := devcontainershell.DevcontainerExec(shell); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	execCmd.Flags().StringVarP(&shell, "shell", "s", "bash", "using shell program")
}
