package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yskszk63/devcontainer-shell"
)

var shell string
var rebuild bool

func exec(cmd *cobra.Command, args []string) error {
	ds, err := devcontainershell.NewDevcontainerShell()
	if err != nil {
		return err
	}

	if err := ds.Exec(rebuild, shell); err != nil {
		return err
	}

	return nil
}

var execCmd = &cobra.Command {
	Use: "exec",
	Short: "execute shell on devcontainer",
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: exec,
}

func setupExecCmd(c *cobra.Command) {
	c.Flags().StringVarP(&shell, "shell", "s", "bash", "using shell program")
	c.Flags().BoolVarP(&rebuild, "rebuild", "b", false, "remove existing container")
}

func init() {
	setupExecCmd(execCmd)
}
