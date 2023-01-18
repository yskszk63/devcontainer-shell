package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yskszk63/devcontainer-shell"
)

var shell string
var noForwardport bool
var rebuild bool

func exec(cmd *cobra.Command, args []string) error {
	ds := new(devcontainershell.DevcontainerShell)
	if rebuild {
		ds.Rebuild = true
	}
	if !noForwardport {
		ds.PortForward = true
	}

	if err := ds.Up(); err != nil {
		return err
	}

	if err := ds.Exec(shell); err != nil {
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
	c.Flags().BoolVarP(&noForwardport, "no-forward", "F", false, "no foward port")
	c.Flags().BoolVarP(&rebuild, "rebuild", "b", false, "remove existing container")
}

func init() {
	setupExecCmd(execCmd)
}
