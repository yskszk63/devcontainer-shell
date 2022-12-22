package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yskszk63/devcontainer-shell"
)

var killCmd = &cobra.Command {
	Use: "kill",
	Short: "kill devcontainer",
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		return devcontainershell.Kill()
	},
}
