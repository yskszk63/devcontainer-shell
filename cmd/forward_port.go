package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

func forwardPort(cmd *cobra.Command, args []string) {
	container := args[0]
	if err := devcontainershell.ForwardPort(container); err != nil {
		log.Fatal(err)
	}
}

var forwardPortCmd = &cobra.Command {
	Use: "forwardport [container id]",
	Short: "port forward",
	Args: cobra.MinimumNArgs(1),
	Run: forwardPort,
}

func init() {
}
