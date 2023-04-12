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
		ds, err := devcontainershell.NewDevcontainerShell()
		if err != nil {
			return err
		}

		return ds.Kill()
	},
}
